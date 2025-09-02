package logger

import "go.uber.org/zap"

type Logger struct {
	logger *zap.SugaredLogger
}

func (l *Logger) Info(args ...interface{}) {
	l.logger.Info(args...)
}

func (l *Logger) Error(args ...interface{}) {
	l.logger.Error(args...)
}

func (l *Logger) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}

func NewLogger() (Logger, error) {
	l := Logger{}
	logger, err := zap.NewDevelopment()
	if err != nil {
		return l, err
	}
	defer logger.Sync() // flushes buffer, if any
	s := logger.Sugar()
	l.logger = s
	return l, nil
}
