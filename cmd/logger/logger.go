package logger

import "go.uber.org/zap"

type ZapLogger struct{}

func (l ZapLogger) New() *zap.SugaredLogger {
	log := zap.NewExample()
	return log.Sugar()
}
