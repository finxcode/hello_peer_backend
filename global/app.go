package global

import (
	"github.com/spf13/viper"
	"webapp_gin/config"
)

type Application struct {
	ConfigViper *viper.Viper
	Config      config.Configuration
}

var app = new(Application)
