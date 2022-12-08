package logic

import (
	"context"
	"fmt"
	"github.com/jinzhu/copier"
	"jakarta/app/pgModel/paymentPgModel"
	"jakarta/common/xerr"

	"jakarta/app/payment/rpc/internal/svc"
	"jakarta/app/payment/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPaymentByFlowNoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetPaymentByFlowNoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPaymentByFlowNoLogic {
	return &GetPaymentByFlowNoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 根据流水号查询流水记录
func (l *GetPaymentByFlowNoLogic) GetPaymentByFlowNo(in *pb.GetPaymentByFlowNoReq) (*pb.GetPaymentByFlowNoResp, error) {
	thirdPayment, err := l.svcCtx.ThirdPaymentFlowModel.FindOne(l.ctx, in.FlowNo)
	if err != nil && err != paymentPgModel.ErrNotFound {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.DbError, fmt.Sprintf("FindOne err:%v , in : %+v", err, in))
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

	return &pb.GetPaymentByFlowNoResp{
		PaymentDetail: &resp,
	}, nil
}
