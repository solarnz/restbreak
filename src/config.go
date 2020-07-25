package restbreak

import (
	"github.com/spf13/viper"
)

func Parse() *RestBreak {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.config/restbreak")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	var config *RestBreak
	err = viper.Unmarshal(&config)
	if err != nil {
		panic(err)
	}

	return config
}
