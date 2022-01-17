package logger

import (
	"fmt"
	"github.com/spf13/viper"
	"go-web-scaffold/settings"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

func Init() (err error) {
	writeSyncer := getLogWriter()

	encoder := getEncoder()

	l := new(zapcore.Level)
	err = l.UnmarshalText([]byte(settings.Config.Level))
	if err != nil {
		return
	}

	core := zapcore.NewCore(encoder, writeSyncer, l)

	logger := zap.New(core, zap.WithCaller(true))

	mode := viper.GetString("app.mode")
	if mode == "dev" || mode == "debug" || mode == "test" {
		logger.WithOptions(zap.Development())
	}

	// replace global default logger without the case that uses logger.logger.info()
	zap.ReplaceGlobals(logger)
	return
}

func Sync() {
	mode := viper.GetString("app.mode")
	if mode == "dev" || mode == "debug" || mode == "test" {
		return
	}
	if err := zap.L().Sync(); err != nil {
		fmt.Printf("Logger Sync failed, err: %v\n", err)
	}
}

func getLogWriter() zapcore.WriteSyncer {
	switch viper.GetString("app.mode") {
	case "prod", "release":
		return zapcore.AddSync(getRollingFileLogger())
	case "dev", "debug", "test":
		fallthrough
	default:
		return zapcore.AddSync(os.Stdout)
	}
}

func getRollingFileLogger() *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   viper.GetString("log.rollingFile.filename"),
		MaxSize:    viper.GetInt("log.rollingFile.maxSize"),
		MaxBackups: viper.GetInt("log.rollingFile.maxBackups"),
		MaxAge:     viper.GetInt("log.rollingFile.maxAge"),
	}
}

func getEncoder() zapcore.Encoder {
	var encoderConfig zapcore.EncoderConfig
	mode := viper.GetString("app.mode")
	switch mode {
	case "prod", "release":
		encoderConfig = zap.NewProductionEncoderConfig()
		encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	case "dev", "debug", "test":
		fallthrough
	default:
		encoderConfig = zap.NewDevelopmentEncoderConfig()
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	encoderConfig.TimeKey = "time"
	encoderConfig.LevelKey = "level"
	encoderConfig.CallerKey = "caller"
	encoderConfig.MessageKey = "msg"
	encoderConfig.StacktraceKey = "stack"
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder

	switch mode {
	case "prod", "release":
		return zapcore.NewJSONEncoder(encoderConfig)
	case "dev", "debug", "test":
		fallthrough
	default:
		return zapcore.NewConsoleEncoder(encoderConfig)
	}
}
