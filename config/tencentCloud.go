package config

type IM struct {
	SdkAppId int    `mapstructure:"sdk_app_id" json:"sdk_app_id" yaml:"sdk_app_id"`
	Key      string `mapstructure:"key" json:"key" yaml:"key"`
	AppName  string `mapstructure:"app_name" json:"app_name" yaml:"app_name"`
	Expiry   int    `mapstructure:"expiry" json:"expiry" yaml:"expiry"`
}
