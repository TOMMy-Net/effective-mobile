package logger

import (
	"io"

	"github.com/sirupsen/logrus"
)

func InitLogger(w io.Writer) *logrus.Logger {
	l := logrus.New()
	l.SetOutput(w)
	l.SetFormatter(&logrus.JSONFormatter{})
	return l
}
