package config

type AppVerConf struct {
	MinAppVer      int64 `json:"MinAppVer"`      // 要求最小版本号
	LatestAppVer   int64 `json:"LatestAppVer"`   // 最新版本号
	StoryTabMaxVer int64 `json:"StoryTabMaxVer"` // story模块最高版本 为了审核开关
}
