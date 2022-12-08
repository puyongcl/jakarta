package logic

import (
	"github.com/zeromicro/go-zero/core/logx"
	"jakarta/app/chat/rpc/pb"
	"testing"
)

func TestProto(t *testing.T) {
	var in pb.SyncChatStateReq
	in.Uid = 1000
	in.Action = 1222
	logx.Errorf("%+v", &in)
	fun1(&in)
}

func fun1(in *pb.SyncChatStateReq) {
	logx.Errorf("%+v", in)
	if in != nil {
		in.ListenerUid = 10002
	}
	logx.Errorf("%+v", in)
}
