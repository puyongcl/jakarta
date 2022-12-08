package main

import (
	"flag"
	"github.com/zeromicro/go-zero/core/logx"
	"jakarta/app/im/mq/internal/config"
	"jakarta/app/im/mq/internal/listen"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
)

var configFile = flag.String("f", "etc/im.yaml", "Specify the config file")

func main() {
	flag.Parse()
	var c config.Config

	conf.MustLoad(*configFile, &c)
	logx.DisableStat()
	// log、prometheus、trace、metricsUrl.
	if err := c.SetUp(); err != nil {
		panic(err)
	}

	serviceGroup := service.NewServiceGroup()
	defer serviceGroup.Stop()

	for _, mq := range listen.Mqs(c) {
		serviceGroup.Add(mq)
	}

	serviceGroup.Start()
}
