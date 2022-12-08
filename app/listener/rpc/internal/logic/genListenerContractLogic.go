package logic

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/lukasjarosch/go-docx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"jakarta/app/pgModel/listenerPgModel"
	"jakarta/common/key/db"
	"jakarta/common/key/listenerkey"
	"jakarta/common/key/rediskey"
	"jakarta/common/key/tencentcloudkey"
	"jakarta/common/third_party/hfbfcash"
	"jakarta/common/third_party/tencentcloud"
	"jakarta/common/tool"
	"jakarta/common/xerr"
	"strings"
	"time"

	"jakarta/app/listener/rpc/internal/svc"
	"jakarta/app/listener/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GenListenerContractLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGenListenerContractLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GenListenerContractLogic {
	return &GenListenerContractLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GenListenerContractLogic) GenListenerContract(req *pb.GenListenerContractReq) (resp *pb.GenListenerContractResp, err error) {
	// 加分布式锁
	rkey := fmt.Sprintf(rediskey.RedisLockUploadCosContract, req.ListenerUid)
	rl := redis.NewRedisLock(l.svcCtx.RedisClient, rkey)
	rl.SetExpire(2)
	b, err := rl.AcquireCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	if !b {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RedisLockFail, "操作太过频繁")
	}
	defer func() {
		_, err2 := rl.ReleaseCtx(l.ctx)
		if err2 != nil {
			logx.WithContext(l.ctx).Errorf("RedisLock %s release err:%+v", rkey, err2)
			return
		}
	}()

	return l.DoGenListenerContract(req)
}

//  生成合同
func (l *GenListenerContractLogic) DoGenListenerContract(req *pb.GenListenerContractReq) (resp *pb.GenListenerContractResp, err error) {
	// 第一次生成合同
	var lcm *listenerPgModel.ListenerContract
	switch req.ContractType {
	case listenerkey.ListenerContractHFBFSZJJHZHB: // 汇服八方数字合作伙伴协议
		// id
		id := fmt.Sprintf(db.DBUidId, req.ListenerUid, req.ContractType)
		remark := fmt.Sprintf("%s-%s", req.ListenerName, req.IdNo)
		// 查询是否生成
		lcm, err = l.svcCtx.ListenerContractModel.FindOne(l.ctx, id)
		if err != nil && err != listenerPgModel.ErrNotFound {
			return
		}
		if lcm != nil { // 已经存在此类型合同
			// 判断身份信息是否改变
			if lcm.Remark != remark {
				// 更新合同
				lcm, err = l.gen10001Contract(id, remark, req)
				if err != nil {
					return
				}

				err = l.svcCtx.ListenerContractModel.UpdateListenerContract(l.ctx, id, lcm.ContractFile, lcm.State, lcm.Remark)
				if err != nil {
					return
				}
			}
			//
			resp = &pb.GenListenerContractResp{
				File:         lcm.ContractFile,
				ContractType: lcm.ContractType,
				UploadState:  lcm.State,
			}
			return
		}
		//
		lcm, err = l.gen10001Contract(id, remark, req)
		if err != nil {
			return
		}

	default:
		return nil, xerr.NewGrpcErrCodeMsg(xerr.AdminGenContractError, "合同类型错误")
	}

	// 生成合同记录
	_, err = l.svcCtx.ListenerContractModel.Insert(l.ctx, lcm)
	if err != nil {
		return
	}

	resp = &pb.GenListenerContractResp{
		File:         lcm.ContractFile,
		ContractType: lcm.ContractType,
		UploadState:  lcm.State,
	}
	return
}

func (l *GenListenerContractLogic) gen10001Contract(id string, remark string, req *pb.GenListenerContractReq) (data *listenerPgModel.ListenerContract, err error) {
	// 生成合同文件
	var templateFile, localFile, cloudFileName string
	var replaceMap docx.PlaceholderMap
	templateFile = l.svcCtx.Config.ContractConfig.ContractTemplate10001
	flag := fmt.Sprintf("%d-%s-%d", req.ListenerUid, req.IdNo, time.Now().Unix())
	fmd5 := tool.Md5ByString(flag)
	localFile = strings.Replace(templateFile, listenerkey.ContractFileNameMark, flag, -1)
	cloudFileName = fmd5 + ".docx"
	replaceMap = docx.PlaceholderMap{
		listenerkey.ContractName:  req.ListenerName,
		listenerkey.ContractName2: req.ListenerName,
		listenerkey.ContractIdNo:  req.IdNo,
		listenerkey.ContractPhone: req.PhoneNumber,
	}

	err = tool.GenNewDocxFromTemplate(templateFile, localFile, replaceMap)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("GenListenerContractLogic gen10001Contract GenNewDocxFromTemplate err:%+v", err)
		err = xerr.NewGrpcErrCodeMsg(xerr.AdminGenContractError, fmt.Sprintf("%+v", err))
		return
	}

	// 初始化
	data = &listenerPgModel.ListenerContract{
		Id:           id,
		ListenerUid:  req.ListenerUid,
		SignTime:     sql.NullTime{},
		ContractFile: "",
		Remark:       remark,
		ContractType: req.ContractType,
		State:        tencentcloud.UploadStateInit,
	}

	// 生成合同 并上传到腾讯云
	data.ContractFile = tencentcloudkey.ContractDir + "/" + cloudFileName
	if l.svcCtx.Config.Mode == service.ProMode {
		err = l.svcCtx.TencentCosClient.Upload(l.ctx, data.ContractFile, localFile)
		if err != nil {
			logx.WithContext(l.ctx).Errorf("GenListenerContractLogic gen10001Contract Upload err:%+v", err)
			return
		}
	}

	data.State = tencentcloud.UploadStateSuccess

	// 同步第三方
	if l.svcCtx.Config.Mode == service.ProMode {
		reqHfbf := hfbfcash.SyncContractReq{
			Name:   req.ListenerName,
			Phone:  req.PhoneNumber,
			Idcard: req.IdNo,
			Remark: "",
			Files: &hfbfcash.SyncContractFile{
				File3: tencentcloudkey.CDNBasePath + data.ContractFile,
			},
		}
		err = l.svcCtx.HfbfCashClient.SyncContract(&reqHfbf)
		if err != nil {
			logx.WithContext(l.ctx).Errorf("GenListenerContractLogic gen10001Contract SyncContract err:%+v, req:%+v", err, reqHfbf)
			return
		}
	}

	return
}
