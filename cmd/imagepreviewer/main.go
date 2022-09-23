package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	lrucache "github.com/alexx1524/go-home-work/hw04_lru_cache"
	internalcache "github.com/alexx1524/go-image-previewer/internal/cache"
	"github.com/alexx1524/go-image-previewer/internal/config"
	"github.com/alexx1524/go-image-previewer/internal/logger"
	internalhttp "github.com/alexx1524/go-image-previewer/internal/server/http"
	"github.com/alexx1524/go-image-previewer/internal/storage"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/image_previewer/config.yaml", "Path to the configuration file")
}

func main() {
	flag.Parse()

	cfg, err := config.NewConfig(configFile)
	if err != nil {
		log.Fatal(err)
	}

	appLog, err := logger.NewLogger(cfg.Log.LogFile, cfg.Log.Level)
	if err != nil {
		log.Fatal(err)
	}

	cache, err := initializeCache(cfg)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := os.Stat(cfg.Storage.ImagesPath); os.IsNotExist(err) {
		err = os.Mkdir(cfg.Storage.ImagesPath, 0o644)
		if err != nil {
			log.Fatal(err)
		}
	}

	fileStorage := storage.NewFileStorage(cfg.Storage.ImagesPath)
	err = fileStorage.RemoveAll()
	if err != nil {
		log.Fatal(err)
	}

	httpServer := internalhttp.NewServer(appLog, cfg, cache, fileStorage)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	appLog.Info("Subscribing for cache removing events")

	go func(ctx context.Context) {
		ch := make(chan lrucache.RemovedItem)
		defer close(ch)

		cache.SetRemoveItemsChan(ch)

		for {
			select {
			case <-ctx.Done():
				appLog.Info("Unsubscribing for cache removing events")
				return
			case deletedItem := <-ch:
				fileName := deletedItem.Value.(string)
				if err := fileStorage.RemoveFile(fileName); err != nil {
					appLog.Error(fmt.Sprintf("File %s - removing error: %s", fileName, err.Error()))
				} else {
					appLog.Debug(fmt.Sprintf("File %s is removed", fileName))
				}
			}
		}
	}(ctx)

	appLog.Info("Starting HTTP Service...")

	go func() {
		if err := httpServer.Start(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
			appLog.Error(fmt.Sprintf("Listen error: %s", err.Error()))
		}
	}()

	<-ctx.Done()

	appLog.Info("Service stopped")
}

func initializeCache(config config.Config) (internalcache.Cache, error) {
	if config.Cache.Mode == "LRUCache" {
		return lrucache.NewCache(config.Cache.LRUCache.ItemsCount), nil
	}
	return nil, errors.New("unsupported cache mode")
}
