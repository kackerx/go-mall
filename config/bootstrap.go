package config

import (
	"os"

	"github.com/spf13/viper"
)

const ConfDir = "config/"

func init() {
	env := os.Getenv("ENV")
	vp := viper.New()

	configFilePath := ConfDir + "app.yaml"
	if env != "" {
		configFilePath = ConfDir + "app." + env + ".yaml"
	}

	vp.SetConfigFile(configFilePath)
	if err := vp.ReadInConfig(); err != nil {
		panic(err)
	}

	if err := vp.Unmarshal(&Conf); err != nil {
		panic(err)
	}
}
