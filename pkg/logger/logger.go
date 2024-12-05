package logger

import (
	"fmt"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"os"
	"path/filepath"

	"github.com/QBC8-Team1/magic-survey/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger methods interface
type Logger interface {
	InitLogger()
	Debug(args ...interface{})
	Debugf(template string, args ...interface{})
	Info(args ...interface{})
	Infof(template string, args ...interface{})
	Warn(args ...interface{})
	Warnf(template string, args ...interface{})
	Error(args ...interface{})
	Errorf(template string, args ...interface{})
	DPanic(args ...interface{})
	DPanicf(template string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(template string, args ...interface{})
}

// Logger
type AppLogger struct {
	cfg         *config.Config
	sugarLogger *zap.SugaredLogger
}

// App Logger constructor
func NewAppLogger(cfg *config.Config) *AppLogger {
	return &AppLogger{cfg: cfg}
}

// For mapping config logger to app logger levels
var loggerLevelMap = map[string]zapcore.Level{
	"debug":  zapcore.DebugLevel,
	"info":   zapcore.InfoLevel,
	"warn":   zapcore.WarnLevel,
	"error":  zapcore.ErrorLevel,
	"dpanic": zapcore.DPanicLevel,
	"panic":  zapcore.PanicLevel,
	"fatal":  zapcore.FatalLevel,
}

func (l *AppLogger) getLoggerLevel(cfg *config.Config) zapcore.Level {
	level, exist := loggerLevelMap[cfg.Logger.Level]
	if !exist {
		return zapcore.DebugLevel
	}

	return level
}

// InitLogger will init logger with config and log rotation
func (l *AppLogger) InitLogger(filePath string) {
	fmt.Println(filePath)
	logLevel := l.getLoggerLevel(l.cfg)

	// Ensure the directory exists
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Fatalf("cannot create log directory: %v", err)
	}

	logWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   filePath,
		MaxSize:    10, // in MB
		MaxBackups: 5,
		MaxAge:     30, // in days
		Compress:   true,
	})

	var encoderCfg zapcore.EncoderConfig
	if l.cfg.Server.Mode == config.Development {
		encoderCfg = zap.NewDevelopmentEncoderConfig()
	} else {
		encoderCfg = zap.NewProductionEncoderConfig()
	}

	encoderCfg.LevelKey = "LEVEL"
	encoderCfg.CallerKey = "CALLER"
	encoderCfg.TimeKey = "TIME"
	encoderCfg.NameKey = "NAME"
	encoderCfg.MessageKey = "MESSAGE"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	encoder := zapcore.NewJSONEncoder(encoderCfg)

	core := zapcore.NewCore(encoder, logWriter, zap.NewAtomicLevelAt(logLevel))
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	l.sugarLogger = logger.Sugar()
	if err := l.sugarLogger.Sync(); err != nil {
		l.sugarLogger.Error(err)
	}
}
func (l *AppLogger) Debug(args ...interface{}) {
	l.sugarLogger.Debug(args...)
}

func (l *AppLogger) Debugf(template string, args ...interface{}) {
	l.sugarLogger.Debugf(template, args...)
}

func (l *AppLogger) Info(args ...interface{}) {
	l.sugarLogger.Info(args...)
}

func (l *AppLogger) Infof(template string, args ...interface{}) {
	l.sugarLogger.Infof(template, args...)
}

func (l *AppLogger) Warn(args ...interface{}) {
	l.sugarLogger.Warn(args...)
}

func (l *AppLogger) Warnf(template string, args ...interface{}) {
	l.sugarLogger.Warnf(template, args...)
}

func (l *AppLogger) Error(args ...interface{}) {
	l.sugarLogger.Error(args...)
}

func (l *AppLogger) Errorf(template string, args ...interface{}) {
	l.sugarLogger.Errorf(template, args...)
}

func (l *AppLogger) DPanic(args ...interface{}) {
	l.sugarLogger.DPanic(args...)
}

func (l *AppLogger) DPanicf(template string, args ...interface{}) {
	l.sugarLogger.DPanicf(template, args...)
}

func (l *AppLogger) Panic(args ...interface{}) {
	l.sugarLogger.Panic(args...)
}

func (l *AppLogger) Panicf(template string, args ...interface{}) {
	l.sugarLogger.Panicf(template, args...)
}

func (l *AppLogger) Fatal(args ...interface{}) {
	l.sugarLogger.Fatal(args...)
}

func (l *AppLogger) Fatalf(template string, args ...interface{}) {
	l.sugarLogger.Fatalf(template, args...)
}
