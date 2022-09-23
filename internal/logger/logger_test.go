package logger

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	filePath := "log"

	removeLogFile := func() error {
		return os.Remove(filePath)
	}

	t.Run("correct level", func(t *testing.T) {
		defer removeLogFile()

		levels := []string{"error", "warn", "info", "debug", "ERROR", "WARN", "INFO", "DEBUG"}
		for _, level := range levels {
			_, err := NewLogger(filePath, level)
			require.NoError(t, err)
		}
	})

	t.Run("incorrect level", func(t *testing.T) {
		_, err := NewLogger(filePath, "incorrect_level")
		require.Error(t, err)
	})

	t.Run("create log file", func(t *testing.T) {
		defer removeLogFile()

		logger, err := NewLogger(filePath, "debug")
		if err != nil {
			log.Fatalln(err)
		}

		logger.Info("info message")
		logger.Debug("debug message")
		logger.Warning("warning message")
		logger.Error("error message")

		_, err = os.ReadFile(filePath)
		if err != nil {
			log.Fatal(err)
		}

		require.NoError(t, err)
	})
}
