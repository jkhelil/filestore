package helper

import (
	"os"
	
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh/terminal"
)

// NewLogger returns a new *logrus.Logger.
func NewLogger(prefix string) *logrus.Logger {
	var err error
	var logger = logrus.New()
	logger.Level, err = logrus.ParseLevel(viper.GetString("log-level"))
	if err != nil {
		logger.Errorf("Couldn't parse log level: %s", viper.GetString("log-level"))
		logger.Level = logrus.InfoLevel
	}
	logger.Formatter = NewPrefixFormatter(prefix)

	return logger
}

// PrefixFormatter wraps the default TextFormatter and prepends a user-defined prefix to all log entries.
type PrefixFormatter struct {
	name string
	f *logrus.TextFormatter
}

// NewPrefixFormatter creates a new PrefixFormatter
func NewPrefixFormatter(name string) *PrefixFormatter {
	return &PrefixFormatter{
		name: name,

		f: &logrus.TextFormatter{},
	}
}

// Format formats the entry using the text formatter then prefixes it.
func (p *PrefixFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	if p.name == "" {
		return p.f.Format(entry)
	}

	// Make a copy of the entry
	prefixed := entry.WithFields(logrus.Fields{})
	prefixed.Level = entry.Level
	prefixed.Message = entry.Message
	prefixed.Time = entry.Time

	if isTerminal(entry.Logger.Out) {
		prefixed.Message = p.name + ": " + prefixed.Message

		return p.f.Format(prefixed)
	}

	prefixed.Data["_name"] = p.name

	return p.f.Format(prefixed)
}

// isTerminal returns a boolean reflecting if output is to a valid terminal. Based on fix https://github.com/sirupsen/logrus/issues/607 to address the removal of logrus.IsTerminal()
func isTerminal(out interface{}) bool {
	lo, ok := out.(*os.File)
	return ok && terminal.IsTerminal(int(lo.Fd()))
}