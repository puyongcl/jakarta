package logic

import (
	"context"
	"fmt"
	"github.com/jinzhu/copier"
	"jakarta/app/pgModel/paymentPgModel"
	"jakarta/common/key/paykey"
	"jakarta/common/xerr"

	"jakarta/app/payment/rpc/internal/svc"
	"jakarta/app/payment/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSuccessPaymentFlowByOrderIdReqLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetSuccessPaymentFlowByOrderIdReqLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSuccessPaymentFlowByOrderIdReqLogic {
	return &GetSuccessPaymentFlowByOrderIdReqLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 根据订单id查询流水记录
func (l *GetSuccessPaymentFlowByOrderIdReqLogic) GetSuccessPaymentFlowByOrderIdReq(in *pb.GetSuccessPaymentFlowByOrderIdReq) (*pb.GetSuccessPaymentFlowByOrderIdResp, error) {
	thirdPayment, err := l.svcCtx.ThirdPaymentFlowModel.FindOneByOrderId(l.ctx, in.OrderId, []int64{paykey.ThirdPaymentPayTradeStateSuccess, paykey.ThirdPaymentPayTradeStateRefundSuccess})
	if err != nil && err != paymentPgModel.ErrNotFound {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.DbError, fmt.Sprintf("FindOneByOrderId orderId:%s err:%+v", in.OrderId, err))
	}

	var resp pb.PaymentDetail
	if thirdPayment != nil {
		_ = copier.Copy(&resp, thirdPayment)
		resp.CreateTime = thirdPayment.CreateTime.Unix()
		resp.UpdateTime = thirdPayment.UpdateTime.Unix()
		if thirdPayment.PayTime.Valid {
			resp.PayTime = thirdPayment.PayTime.Time.Unix()
		}
	}

	return &pb.GetSuccessPaymentFlowByOrderIdResp{
		PaymentDetail: &resp,
	}, nil
}
