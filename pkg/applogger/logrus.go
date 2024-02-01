package applogger

import (
	"sync"

	"github.com/sangianpatrick/tm-user/config"
	"github.com/sangianpatrick/tm-user/pkg/logrushook"
	"github.com/sirupsen/logrus"
	"github.com/uptrace/opentelemetry-go-extra/otellogrus"
)

var (
	logrusLogger   *logrus.Logger
	logrusSyncOnce sync.Once
)

func constructLogrus() *logrus.Logger {
	cfg := config.Get()

	logger := logrus.New()
	logger.SetReportCaller(true)
	logger.SetFormatter(cfg.Logrus.Formatter)
	logger.AddHook(otellogrus.NewHook(otellogrus.WithLevels(logrus.AllLevels...)))
	logger.AddHook(logrushook.NewTraceIDLoggerHook())
	logger.AddHook(logrushook.NewStdoutLoggerHook(logrus.New(), cfg.Logrus.Formatter))

	return logger
}

func GetLogrus() *logrus.Logger {
	logrusSyncOnce.Do(func() {
		logrusLogger = constructLogrus()
	})

	return logrusLogger
}
