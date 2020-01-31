package helper

import (
	"testing"
)

func TestLoggerDefaultsToInfo(t *testing.T) {
	l := logx.New("")

	require.Equal(t, logrus.InfoLevel, l.Level)
}

