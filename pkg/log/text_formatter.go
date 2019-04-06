package log

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

var TextFormatter = &textFormatter{}

type textFormatter struct {
	logrus.TextFormatter
	Prefix string
}

func (f *textFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	if f.Prefix != "" {
		entry.Message = fmt.Sprintf("[%s] %s", f.Prefix, entry.Message)
	}
	return f.TextFormatter.Format(entry)
}
