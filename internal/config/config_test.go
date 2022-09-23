package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	t.Run("If config file doesn't exist returns error", func(t *testing.T) {
		_, err := NewConfig("wrong file path")

		require.Error(t, err)
	})

	t.Run("Correct parsing all parameters", func(t *testing.T) {
		cfg, err := NewConfig("../../configs/tests/config.yaml")

		require.NoError(t, err)

		require.Equal(t, cfg.HTTPServer.Address, "localhost")
		require.Equal(t, cfg.HTTPServer.Port, 1234)
		require.Equal(t, cfg.HTTPServer.ReadHeaderTimeout, 10)
		require.Equal(t, cfg.Storage.ImagesPath, "images")
		require.Equal(t, cfg.Cache.Mode, "LRUCache")
		require.Equal(t, cfg.Cache.LRUCache.ItemsCount, 1000)
		require.Equal(t, cfg.Log.LogFile, "log.txt")
		require.Equal(t, cfg.Log.Level, "debug")
	})
}
