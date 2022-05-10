package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(name string) (*zap.SugaredLogger, error) {
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{"stdout"}
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.DisableStacktrace = true
	cfg.InitialFields = map[string]interface{}{
		"service": name,
	}

	log, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	return log.Sugar(), nil
}
