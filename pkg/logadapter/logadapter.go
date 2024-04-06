package logadapter

import (
	"github.com/go-kit/log"
	"github.com/sirupsen/logrus"
)

type logrusAdapter struct {
	*logrus.Logger
}

func (l logrusAdapter) Log(keyvals ...interface{}) error {
	if len(keyvals)%2 != 0 {
		keyvals = append(keyvals, "Log called with odd number of keyvals")
	}

	fields := logrus.Fields{}
	for i := 0; i < len(keyvals); i += 2 {
		key, ok := keyvals[i].(string)
		if !ok {
			l.Logger.Error("Key is not a string", keyvals[i])
			continue
		}
		fields[key] = keyvals[i+1]
	}

	l.Logger.WithFields(fields).Info()
	return nil
}

func NewLogrusAdapter(logger *logrus.Logger) log.Logger {
	return &logrusAdapter{Logger: logger}
}
