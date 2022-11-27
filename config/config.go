package config

type Configuration struct {
	App      App      `mapstructure:"app" json:"app" yaml:"app"`
	Log      Log      `mapstructure:"log" json:"log" yaml:"log"`
	Database Database `mapstructure:"database" json:"database" yaml:"database"`
	Jwt      Jwt      `mapstructure:"jwt" json:"jwt" yaml:"jwt"`
	Redis    Redis    `mapstructure:"redis" json:"redis" yaml:"redis"`
	Wechat   Wechat   `mapstructure:"wechat" json:"wechat" yaml:"wechat"`
	IM       IM       `mapstructure:"im" json:"im" yaml:"im"`
	Secret   Secret   `mapstructure:"secret" json:"secret" yaml:"secret"`
}
