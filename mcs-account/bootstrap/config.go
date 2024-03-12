package bootstrap

import (
	"github.com/lengocson131002/go-clean-core/config"
	"github.com/lengocson131002/go-clean-core/config/viper"
)

func GetConfigure() config.Configure {
	var file viper.ConfigFile = ".env"
	cfg, err := viper.NewViperConfig(&file)
	if err != nil {
		panic(err)
	}
	return cfg
}
