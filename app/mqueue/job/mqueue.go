package main

import (
	"context"
	"flag"
	"github.com/zeromicro/go-zero/core/logx"
	"jakarta/app/mqueue/job/internal/logic"
	"jakarta/app/mqueue/job/internal/svc"
	"os"

	"jakarta/app/mqueue/job/internal/config"

	"github.com/zeromicro/go-zero/core/conf"
)

var configFile = flag.String("f", "etc/mqueue.yaml", "Specify the config file")

func main() {
	flag.Parse()
	var c config.Config

	conf.MustLoad(*configFile, &c)
	logx.DisableStat()
	// log、prometheus、trace、metricsUrl
	if err := c.SetUp(); err != nil {
		panic(err)
	}

	//logx.DisableStat()

	svcContext := svc.NewServiceContext(c)
	ctx := context.Background()
	cronJob := logic.NewCronJob(ctx, svcContext)
	mux := cronJob.Register()

	if err := svcContext.AsynqServer.Run(mux); err != nil {
		logx.WithContext(ctx).Errorf("!!!CronJobErr!!! run err:%+v", err)
		os.Exit(1)
	}
}
