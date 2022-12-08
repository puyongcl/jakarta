package third

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
	"jakarta/app/payment/api/internal/svc"
	"jakarta/app/payment/api/internal/types"
	"jakarta/common/kqueue"
	"sort"
	"strconv"
	"strings"
)

type HfbfCashCallbackLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewHfbfCashCallbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HfbfCashCallbackLogic {
	return &HfbfCashCallbackLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HfbfCashCallbackLogic) HfbfCashCallback(req *types.HFBFCashCallbackReq) (resp *types.HFBFCashCallbackResp, err error) {
	//logx.WithContext(l.ctx).Infof("HfbfCashCallbackLogic HfbfCashCallback ts:%d, params:%+v", req.TimeStamp, req.Params)
	kqMsg := kqueue.UpdateCashStatusMessage{}
	_ = copier.Copy(&kqMsg, req.Params)
	returnCode := 200
	msg := "sucess"

	l.Verify(req)

	err = l.pushKq(&kqMsg)
	if err != nil {
		returnCode = 400
	}

	resp = &types.HFBFCashCallbackResp{
		Code:    int64(returnCode),
		Message: msg,
	}
	return resp, err
}

func (l *HfbfCashCallbackLogic) Verify(req *types.HFBFCashCallbackReq) {
	var cc map[string]interface{}
	b, _ := json.Marshal(req.Params)
	_ = json.Unmarshal(b, &cc)
	st := getValues(cc, "")
	timestamp := float64(req.TimeStamp)
	timestampStr := strconv.FormatFloat(timestamp, 'f', 0, 64)
	sign := Md5(timestampStr + st + l.svcCtx.Config.HfbfCashConf.AppSecret)
	sign = strings.ToUpper(sign)
	if req.Sign != sign {
		logx.WithContext(l.ctx).Errorf("HfbfCashCallbackLogic Verify Sign calc sign:%s, req sign:%s", sign, req.Sign)
	}
	return
}

func (l *HfbfCashCallbackLogic) pushKq(msg *kqueue.UpdateCashStatusMessage) error {
	buf, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return l.svcCtx.KqueueUpdateCashStatusClient.Push(string(buf))
}

func getValues(cc map[string]interface{}, str string) string {
	var keys []string
	for k := range cc {
		keys = append(keys, k)
	}
	//按字典升序排列
	sort.Strings(keys)
	for _, k := range keys {
		if value, ok := cc[k].([]interface{}); ok {
			for _, v := range value {
				if new_map, ok2 := v.(map[string]interface{}); ok2 {
					str = getValues(new_map, str)
				} else {
					str += StrHandle(v)
				}
			}
		} else {
			str += StrHandle(cc[k])
		}
	}
	return str
}

func StrHandle(data interface{}) string {
	var str string
	switch data.(type) {
	case string:
		ft := data.(string)
		str = ft
	case float64:
		ft := data.(float64)
		str = strconv.FormatFloat(ft, 'f', -1, 64)
	case int:
		it := data.(int)
		str = strconv.Itoa(it)
	case bool:
		bt := data.(bool)
		if bt {
			str = "True"
		} else {
			str = "False"
		}
	}
	str = strings.Replace(str, " ", "", -1)
	return str
}

// 返回一个32位md5加密后的字符串
func Md5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
