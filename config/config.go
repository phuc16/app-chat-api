package config

import (
	"app/build"
	"app/pkg/logger"
	"app/pkg/trace"
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type loggerConfig struct {
	Mode       string `yaml:"mode" env:"LOG_MODE"`
	Encoding   string `yaml:"encoding" env:"LOG_ENCODING"`
	Level      string `yaml:"level" env:"LOG_LEVEL"`
	LogFile    string `yaml:"log_file" env:"LOG_FILE"`
	StackTrace bool   `yaml:"stack_trace" env:"LOG_STACK_TRACE"`
}

func (c loggerConfig) ToLoggerConfig() logger.Config {
	return logger.Config{
		Mode:     c.Mode,
		Encoding: c.Encoding,
		Level:    c.Level,
		LogFile:  c.LogFile,
	}
}

type otelConfig struct {
	TraceProvider struct {
		Enable           bool   `yaml:"enable" env:"TRACE_PROVIDER_ENABLE"`
		OtlpHttpEndpoint string `yaml:"endpoint" env:"TRACE_PROVIDER_OTLP_HTTP_ENDPOINT"`
		OtlpHttpInsecure bool   `yaml:"insecure" env:"TRACE_PROVIDER_INSECURE"`
	} `yaml:"trace_provider"`
	MetricProvider struct {
		Enable bool `yaml:"enable" env:"METRIC_PROVIDER_ENABLE"`
	} `yaml:"metric_provider"`
}

func (c otelConfig) ToTraceConfig() trace.OTelConfig {
	return trace.OTelConfig{
		ServiceName:    build.AppName,
		ServiceVersion: build.Version,
		TraceProvider: struct {
			Enable           bool
			OtlpHttpEndpoint string
			OtlpHttpInsecure bool
		}{
			Enable:           c.TraceProvider.Enable,
			OtlpHttpEndpoint: c.TraceProvider.OtlpHttpEndpoint,
			OtlpHttpInsecure: c.TraceProvider.OtlpHttpInsecure,
		},
		MetricProvider: struct {
			Enable bool
		}{
			Enable: c.MetricProvider.Enable,
		},
	}
}

type hTTPConfig struct {
	Host                string   `yaml:"host" env:"HTTP_HOST"`
	Port                int      `yaml:"port" env:"HTTP_PORT"`
	Origin              string   `yaml:"origin" env:"HTTP_ORIGIN"`
	AllowOrigins        []string `yaml:"allow_origins" env:"HTTP_ALLOW_ORIGINS"`
	Secret              string   `yaml:"secret" env:"HTTP_SECRET"`
	EnableSSL           bool     `yaml:"enable_ssl" env:"HTTP_ENABLE_SSL"`
	CertFile            string   `yaml:"cert_file" env:"HTTP_TLS_CERT_FILE"`
	KeyFile             string   `yaml:"key_file" env:"HTTP_TLS_KEY_FILE"`
	AccessTokenDuration int      `yaml:"access_token_duration" env:"HTTP_ACCESS_TOKEN_DURATION"`
	IsProduction        bool     `yaml:"is_production" env:"IS_PRODUCTION"`
}

func (c hTTPConfig) Addr() string {
	addr := fmt.Sprintf("%s:%d", c.Host, c.Port)
	return addr
}

func (c hTTPConfig) FullAddr() string {
	if c.EnableSSL {
		return fmt.Sprintf("https://%s", c.Addr())
	}
	return fmt.Sprintf("http://%s", c.Addr())
}

type dbConfig struct {
	URI    string `yaml:"uri" env:"DB_URI"`
	DBName string `yaml:"db_name" env:"DB_NAME"`
}

type mailConfig struct {
	Host     string `yaml:"host" env:"MAIL_HOST"`
	Port     int    `yaml:"port" env:"MAIL_PORT"`
	User     string `yaml:"user" env:"MAIL_USER"`
	Password string `yaml:"password" env:"MAIL_PASSWORD"`
}

type allConfig struct {
	Logger loggerConfig `yaml:"logger"`
	OTel   otelConfig   `yaml:"otel"`
	HTTP   hTTPConfig   `yaml:"http"`
	DB     dbConfig     `yaml:"db"`
	Mail   mailConfig   `yaml:"mail"`
}

var Cfg allConfig

func Load() error {
	configFile := "config.yaml"
	if file := os.Getenv("CONFIG_FILE"); file != "" {
		configFile = file
	}
	err := cleanenv.ReadConfig(configFile, &Cfg)
	if err != nil {
		return err
	}
	return nil
}
