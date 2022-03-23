package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"simpleGoJWT/utils"
	"sync"
)

type Config struct {
	IsDebug *bool `yaml:"is_debug" env-required:"true"`
	Listen  struct {
		Host string `yaml:"host" env-default:"127.0.0.1"`
		Port string `yaml:"port" env-default:"3000"`
	} `yaml:"listen"`
	MongoDB struct {
		Host     string `json:"host" env-required:"true"`
		Port     string `json:"port" env-required:"true"`
		Database string `json:"database" env-required:"true"`
		AuthDB   string `json:"auth_db"`
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"mongodb"`
	JWT struct {
		SecretKeyJWT string `yaml:"secret_key_jwt" env-required:"true"`
	}
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		logger := utils.GetLogger()
		logger.Info("read application configuration")
		instance = &Config{}
		if err := cleanenv.ReadConfig("config.yml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger.Info(help)
			logger.Fatal(err)
		}
	})
	return instance
}
