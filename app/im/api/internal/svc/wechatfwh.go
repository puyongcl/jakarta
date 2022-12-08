package svc

import (
	"context"
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/officialaccount"
	offConfig "github.com/silenceper/wechat/v2/officialaccount/config"
	"github.com/zeromicro/go-zero/core/service"
	"jakarta/app/im/api/internal/config"
)

//InitWechat 获取wechat实例
//在这里已经设置了全局cache，则在具体获取公众号/小程序等操作实例之后无需再设置，设置即覆盖
func InitWechat(c *config.Config) (*wechat.Wechat, *officialaccount.OfficialAccount) {
	if c.Mode != service.ProMode {
		return nil, nil
	}
	wc := wechat.NewWechat()
	redisOpts := &cache.RedisOpts{
		Host:     c.RedisCache.Host,
		Password: c.RedisCache.Pass,
	}
	redisCache := cache.NewRedis(context.Background(), redisOpts)
	wc.SetCache(redisCache)
	offCfg := &offConfig.Config{
		AppID:          c.WxFwhConf.AppID,
		AppSecret:      c.WxFwhConf.AppSecret,
		Token:          c.WxFwhConf.Token,
		EncodingAESKey: c.WxFwhConf.EncodingAESKey,
	}
	off := wc.GetOfficialAccount(offCfg)
	return wc, off
}
