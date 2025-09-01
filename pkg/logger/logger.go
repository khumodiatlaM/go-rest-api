package logger

import "go.uber.org/zap"

var sugar *zap.SugaredLogger

func NewLogger() (*zap.SugaredLogger, error) {
	if sugar != nil {
		logger, err := zap.NewDevelopment()
		if err != nil {
			return nil, err
		}
		sugar = logger.Sugar()
	}
	return sugar, nil
}

func Sync() {
	if sugar != nil {
		_ = sugar.Sync()
	}
}
