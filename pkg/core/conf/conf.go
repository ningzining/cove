package conf

import (
	"log"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

func MustLoad(file string, cfg any) {
	if err := Load(file, cfg); err != nil {
		log.Fatalf("load config file fail[%s], %v", file, err)
	}
}

func Load(file string, cfg any) error {
	viper.SetConfigFile(file)
	if err := viper.ReadInConfig(); err != nil {
		return errors.Wrap(err, "failed to read config file")
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return errors.Wrap(err, "failed to unmarshal config")
	}
	return nil
}
