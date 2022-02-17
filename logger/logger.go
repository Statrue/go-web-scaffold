package logger

import (
	"fmt"
	"go-web-scaffold/settings"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

func Init(cfg *settings.LoggerConfig) (err error) {
	mode := settings.Conf.Mode
	writeSyncer := getLogWriter(cfg, mode)

	encoder := getEncoder(mode)

	l := new(zapcore.Level)
	err = l.UnmarshalText([]byte(cfg.Level))
	if err != nil {
		return
	}

	core := zapcore.NewCore(encoder, writeSyncer, l)

	logger := zap.New(core, zap.WithCaller(true))

	if mode == "dev" || mode == "debug" || mode == "test" {
		logger.WithOptions(zap.Development())
	}

	// replace global default logger without the case that uses logger.logger.info()
	zap.ReplaceGlobals(logger)
	return
}

func Sync() {
	mode := settings.Conf.Mode
	if mode == "dev" || mode == "debug" || mode == "test" {
		return
	}
	if err := zap.L().Sync(); err != nil {
		fmt.Printf("Logger Sync failed, err: %v\n", err)
	}
}

func getLogWriter(cfg *settings.LoggerConfig, mode string) zapcore.WriteSyncer {
	switch mode {
	case "prod", "release":
		return zapcore.AddSync(getRollingFileLogger(cfg))
	case "dev", "debug", "test":
		fallthrough
	default:
		return zapcore.AddSync(os.Stdout)
	}
}

func getRollingFileLogger(cfg *settings.LoggerConfig) *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   cfg.RollingFileConfig.Filename,
		MaxSize:    cfg.RollingFileConfig.MaxSize,
		MaxBackups: cfg.RollingFileConfig.MaxBackups,
		MaxAge:     cfg.RollingFileConfig.MaxAge,
	}
}

func getEncoder(mode string) zapcore.Encoder {
	var encoderConfig zapcore.EncoderConfig
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
