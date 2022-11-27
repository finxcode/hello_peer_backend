package config

type Secret struct {
	Path        string `mapstructure:"path" json:"path" yaml:"path"`
	AccessToken string `mapstructure:"access_token" json:"access_token" yaml:"access_token"`
}
