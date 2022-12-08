package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"jakarta/app/pgModel/userPgModel"
	"jakarta/app/usercenter/rpc/internal/svc"
	"jakarta/app/usercenter/rpc/pb"
	"jakarta/common/key/userkey"
	"jakarta/common/uniqueid"
	"jakarta/common/xerr"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProcessReportUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProcessReportUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProcessReportUserLogic {
	return &ProcessReportUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  处理上报用户
func (l *ProcessReportUserLogic) ProcessReportUser(in *pb.AdminProcessReportUserReq) (*pb.AdminProcessReportUserResp, error) {
	// 处理
	var addDays int
	switch in.Action {
	case userkey.ReportStateCreated: // 用户上报
	case userkey.ReportStateBlock2Days: // 封号2天
		addDays = 2
	case userkey.ReportStateBlock7Days: // 封号7天
		addDays = 7
	case userkey.ReportStateBlockForever: // 封号永久
		addDays = 1000000
	case userkey.ReportStateNotTrue: // 不符合事实
	case userkey.ReportStateCancel, userkey.ReportStateAdminBan2Days, userkey.ReportStateAdminBan7Days,
		userkey.ReportStateAdminBanForever: // 管理员直接封号、取消封号
		return l.adminBanUser(in)

	default:
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "参数错误")
	}

	// 更新封号
	if addDays > 0 {
		var tip []string
		for idx := 0; idx < len(in.ReportTag); idx++ {
			tip = append(tip, userkey.ReportTagText[idx])
		}
		if len(tip) <= 0 {
			tip = append(tip, in.ReportContent)
		}

		err := l.svcCtx.UserAuthModel.UpdateBanAccount(l.ctx, in.TargetUid, time.Now().AddDate(0, 0, addDays), strings.Join(tip, "，"))
		if err != nil {
			return nil, err
		}
	}

	// 更新状态
	err := l.svcCtx.UserReportModel.UpdateState(l.ctx, in.Id, in.Action, in.Remark)
	if err != nil {
		return nil, err
	}

	// 发送通知
	return &pb.AdminProcessReportUserResp{}, nil
}

func (l *ProcessReportUserLogic) adminBanUser(in *pb.AdminProcessReportUserReq) (*pb.AdminProcessReportUserResp, error) {
	switch in.Action {
	case userkey.ReportStateAdminBan2Days, userkey.ReportStateAdminBan7Days, userkey.ReportStateAdminBanForever: // 管理员封号
		var addDays int
		switch in.Action {
		case userkey.ReportStateAdminBan2Days: // 封号2天
			addDays = 2
		case userkey.ReportStateAdminBan7Days: // 封号7天
			addDays = 7
		case userkey.ReportStateAdminBanForever: // 封号永久
			addDays = 1000000
		}
		err := l.svcCtx.UserAuthModel.UpdateBanAccount(l.ctx, in.TargetUid, time.Now().AddDate(0, 0, addDays), in.Remark)
		if err != nil {
			return nil, err
		}
		data := new(userPgModel.UserReport)
		_ = copier.Copy(data, in)
		data.State = in.Action
		data.Id = uniqueid.GenDataId()
		_, err = l.svcCtx.UserReportModel.Insert(l.ctx, data)
		if err != nil {
			return nil, err
		}
	case userkey.ReportStateCancel: // 管理员解封
		err := l.svcCtx.UserAuthModel.UpdateBanAccount(l.ctx, in.TargetUid, time.Now(), in.Remark)
		if err != nil {
			return nil, err
		}
	default:
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "类型错误")

	}
	return &pb.AdminProcessReportUserResp{}, nil
}
