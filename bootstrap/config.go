package bootstrap

import "github.com/spf13/viper"

func InitConfig() *viper.Viper {
	config := "config.yaml"
	v := viper.New()
	v.SetConfigFile(config)
	return v
}
