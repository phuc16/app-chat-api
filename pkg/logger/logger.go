package logger

import (
	"context"
	"fmt"
	"os"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Config struct {
	Mode     string
	Encoding string
	Level    string
	LogFile  string
}

var instance *logger

func InitLogger(cfg Config) {
	instance = NewLogger(cfg)
}

func For(ctx context.Context) *logger {
	return instance.For(ctx)
}

// Logger
type logger struct {
	cfg        Config
	logger     *zap.Logger
	spanFields []zapcore.Field
}

// Logger constructor
func NewLogger(cfg Config) *logger {
	logger := &logger{cfg: cfg}
	logger.InitLogger()
	return logger
}

// For mapping config logger to email_service logger levels
var loggerLevelMap = map[string]zapcore.Level{
	"debug":  zapcore.DebugLevel,
	"info":   zapcore.InfoLevel,
	"warn":   zapcore.WarnLevel,
	"error":  zapcore.ErrorLevel,
	"dpanic": zapcore.DPanicLevel,
	"panic":  zapcore.PanicLevel,
	"fatal":  zapcore.FatalLevel,
}

func (l *logger) getLoggerLevel(cfg Config) zapcore.Level {
	level, exist := loggerLevelMap[cfg.Level]
	if !exist {
		return zapcore.DebugLevel
	}

	return level
}

func (l *logger) InitLogger() {
	logLevel := l.getLoggerLevel(l.cfg)

	var encoderCfg zapcore.EncoderConfig
	if l.cfg.Mode == "production" {
		encoderCfg = zap.NewProductionEncoderConfig()
	} else {
		encoderCfg = zap.NewDevelopmentEncoderConfig()
	}

	var encoder zapcore.Encoder
	encoderCfg.LevelKey = "LEVEL"
	encoderCfg.CallerKey = "CALLER"
	encoderCfg.TimeKey = "TIME"
	encoderCfg.NameKey = "NAME"
	encoderCfg.MessageKey = "MESSAGE"
	encoderCfg.FunctionKey = "FUNC"
	encoderCfg.EncodeDuration = zapcore.NanosDurationEncoder
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	if l.cfg.Encoding == "console" {
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	}

	ws := []zapcore.WriteSyncer{}
	ws = append(ws, zapcore.AddSync(os.Stdout))
	if l.cfg.LogFile != "" {
		ws = append(ws, zapcore.AddSync(&lumberjack.Logger{
			Filename:   l.cfg.LogFile,
			MaxSize:    1, // megabytes
			MaxBackups: 10,
		}))
	}
	writeSyncer := zapcore.NewMultiWriteSyncer(ws...)

	core := zapcore.NewCore(encoder, writeSyncer, zap.NewAtomicLevelAt(logLevel))
	logger := zap.New(core, zap.AddStacktrace(zapcore.FatalLevel), zap.AddCaller(), zap.AddCallerSkip(1))
	l.logger = logger
}

// Logger methods

func (l *logger) For(ctx context.Context) *logger {
	span := trace.SpanFromContext(ctx).SpanContext()
	spanFields := []zapcore.Field{}
	if span.HasTraceID() {
		spanFields = append(spanFields, []zapcore.Field{
			zap.String("trace_id", span.TraceID().String()),
			zap.String("span_id", span.SpanID().String()),
		}...)
	}
	return &logger{
		logger:     l.logger,
		spanFields: spanFields,
	}
}

func (l *logger) Write(p []byte) (n int, err error) {
	l.logger.Info(string(p), l.spanFields...)
	return len(p), nil
}

func (l *logger) Append(field ...zapcore.Field) *logger {
	l.spanFields = append(l.spanFields, field...)
	return l
}

func (l *logger) Printf(template string, args ...interface{}) {
	l.logger.Info(getMessage(template, args), l.spanFields...)
}

func (l *logger) Debug(msg string) {
	l.logger.Debug(msg, l.spanFields...)
}

func (l *logger) Debugf(template string, args ...interface{}) {
	l.logger.Debug(getMessage(template, args), l.spanFields...)
}

func (l *logger) Info(msg string) {
	l.logger.Info(msg, l.spanFields...)
}

func (l *logger) Infof(template string, args ...interface{}) {
	l.logger.Info(getMessage(template, args), l.spanFields...)
}

func (l *logger) Warn(msg string) {
	l.logger.Warn(msg, l.spanFields...)
}

func (l *logger) Warnf(template string, args ...interface{}) {
	l.logger.Warn(getMessage(template, args), l.spanFields...)
}

func (l *logger) Error(msg string) {
	l.logger.Error(msg, l.spanFields...)
}

func (l *logger) Errorf(template string, args ...interface{}) {
	l.logger.Error(getMessage(template, args), l.spanFields...)
}

func (l *logger) DPanic(msg string) {
	l.logger.DPanic(msg, l.spanFields...)
}

func (l *logger) DPanicf(template string, args ...interface{}) {
	l.logger.DPanic(getMessage(template, args), l.spanFields...)
}

func (l *logger) Panic(msg string) {
	l.logger.Panic(msg, l.spanFields...)
}

func (l *logger) Panicf(template string, args ...interface{}) {
	l.logger.Panic(getMessage(template, args), l.spanFields...)
}

func (l *logger) Fatal(msg string) {
	l.logger.Fatal(msg, l.spanFields...)
}

func (l *logger) Fatalf(template string, args ...interface{}) {
	l.logger.Fatal(getMessage(template, args), l.spanFields...)
}

// getMessage format with Sprint, Sprintf, or neither.
func getMessage(template string, fmtArgs []interface{}) string {
	if len(fmtArgs) == 0 {
		return template
	}

	if template != "" {
		return fmt.Sprintf(template, fmtArgs...)
	}

	if len(fmtArgs) == 1 {
		if str, ok := fmtArgs[0].(string); ok {
			return str
		}
	}
	return fmt.Sprint(fmtArgs...)
}
