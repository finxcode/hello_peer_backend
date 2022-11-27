package config

type Wechat struct {
	ApiKey          string `mapstructure:"api_key" json:"api_key" yaml:"api_key"`
	ApiSecret       string `mapstructure:"api_secret" json:"api_secret" yaml:"api_secret"`
	AccessGrantType string `mapstructure:"access_grant_type" json:"access_grant_type" yaml:"access_grant_type"`
}
