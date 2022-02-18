package settings

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var Conf = new(GlobalConfig)

type GlobalConfig struct {
	*AppConfig    `mapstructure:"app"`
	*LoggerConfig `mapstructure:"log"`
	*MySQLConfig  `mapstructure:"mysql"`
	*RedisConfig  `mapstructure:"redis"`
}

type AppConfig struct {
	Name    string
	Mode    string
	Version string
	Port    int
}

type LoggerConfig struct {
	Level              string
	*RollingFileConfig `mapstructure:"rollingFile"`
}

type RollingFileConfig struct {
	Filename   string
	MaxSize    int
	MaxAge     int
	MaxBackups int
}

type MySQLConfig struct {
	Host             string
	Port             int
	User             string
	Password         string
	Schema           string
	Params           string
	MaxConn, MaxIdle int
}

type RedisConfig struct {
	Host     string
	Port     int
	DB       int `mapstructure:"db"`
	PoolSize int
	Password string
}

func Init() (err error) {
	// Specify conf file
	viper.SetConfigName("config")
	viper.AddConfigPath("./conf/")

	if err = viper.ReadInConfig(); err != nil {
		return
	}

	// Hot-loading
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("Config file modified.")
		if err = viper.Unmarshal(Conf); err != nil {
			return
		}
	})
	viper.WatchConfig()

	if err = viper.Unmarshal(Conf); err != nil {
		return
	}

	return
}
