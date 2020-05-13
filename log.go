package geekhub

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

func Debugf(format string, args ...interface{}) {
	logger.Debugf(fmt.Sprintf("[offset=%d, page=%d] ", curOffset, curPostsPage) + format, args...)
}

func Warnf(format string, args ...interface{}) {
	logger.Warnf(fmt.Sprintf("[offset=%d, page=%d] ", curOffset, curPostsPage) + format, args...)
}
