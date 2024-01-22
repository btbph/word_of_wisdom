package logger

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log/slog"
	"testing"
)

func TestParseLogLevel(t *testing.T) {
	tt := []struct {
		name     string
		logLevel string
		expected slog.Level
	}{
		{
			name:     "debug log level",
			logLevel: "debug",
			expected: slog.LevelDebug,
		},
		{
			name:     "info log level",
			logLevel: "info",
			expected: slog.LevelInfo,
		},
		{
			name:     "warn log level",
			logLevel: "warn",
			expected: slog.LevelWarn,
		},
		{
			name:     "error log level",
			logLevel: "error",
			expected: slog.LevelError,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result, err := ParseLogLevel(tc.logLevel)

			require.NoError(t, err)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestParseLogLevel_wrongLogLevel(t *testing.T) {
	logLevel := "badLogLevel"

	result, err := ParseLogLevel(logLevel)

	require.Error(t, err)
	require.Empty(t, result)
	assert.EqualError(t, err, "unkown log level")
}
