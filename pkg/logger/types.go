package logger

import "github.com/sirupsen/logrus"

// Logger enforces specific log message formats
type Logger struct {
	*logrus.Logger
}
