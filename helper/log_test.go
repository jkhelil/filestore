package helper

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestNewLogger(t *testing.T) {
	l := NewLogger("filestore")

	require.Equal(t, logrus.InfoLevel, l.Level)
}

