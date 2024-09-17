package logger

import "go.uber.org/zap"

type Logger struct {
	*zap.Logger
}

func Init(development bool) (*Logger, error) {
	var logger *zap.Logger
	var err error

	if development {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}

	if err != nil {
		return nil, err
	}

	return &Logger{
		Logger: logger,
	}, nil
}

func (l *Logger) Sync() {
	l.Logger.Sync()
}
