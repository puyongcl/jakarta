package config

type CosConf struct {
	SecretId  string `json:"SecretId"`
	SecretKey string `json:"SecretKey"`
	AppId     string `json:"AppId"`
	Bucket    string `json:"Bucket"`
	Region    string `json:"Region"`
	Expire    int64  `json:"Expire"`
}
