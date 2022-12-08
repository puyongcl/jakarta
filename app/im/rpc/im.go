package main

import (
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"

	"jakarta/app/im/rpc/internal/config"
	"jakarta/app/im/rpc/internal/server"
	"jakarta/app/im/rpc/internal/svc"
	"jakarta/app/im/rpc/pb"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/im.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	logx.DisableStat()
	ctx := svc.NewServiceContext(c)
	svr := server.NewImServer(ctx)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		pb.RegisterImServer(grpcServer, svr)

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	//
	zrpc.DontLogContentForMethod("/pb.im/addImMsgLog")

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
