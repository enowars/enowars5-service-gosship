package logger

import "github.com/sirupsen/logrus"

func New(level ...logrus.Level) *logrus.Logger {
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	if len(level) > 0 {
		log.SetLevel(level[0])
	}
	return log
}
