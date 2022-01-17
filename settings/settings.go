package settings

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var Config GlobalConfig

type GlobalConfig struct {
	AppConfig    `mapstructure:"app"`
	LoggerConfig `mapstructure:"log"`
	MySQLConfig  `mapstructure:"mysql"`
	RedisConfig  `mapstructure:"redis"`
}

type AppConfig struct {
	Name string
	Mode string
	Port int
}

type LoggerConfig struct {
	Level             string
	RollingFileConfig `mapstructure:"rollingFile"`
}

type RollingFileConfig struct {
	Filename   string
	MaxSize    int
	MaxAge     int
	MaxBackups int
}

type MySQLConfig struct {
	Host             string
	Port             string
	User             string
	Password         string
	Schema           string
	Params           string
	MaxConn, MaxIdle int
}

type RedisConfig struct {
	Host string
	Port int
	DB   int
}

func Init() (err error) {
	// Specify conf file's name and type. (or use SetConfigFile() for briefness)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./conf/")

	if err = viper.ReadInConfig(); err != nil {
		return
	}

	// Hot-loading
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("Config file modified.")
	})
	viper.WatchConfig()

	if err = viper.Unmarshal(&Config); err != nil {
		return
	}

	return
}
