package http

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path"
	"strconv"

	lrucache "github.com/alexx1524/go-home-work/hw04_lru_cache"
	"github.com/alexx1524/go-image-previewer/internal/imageprocessor"
	guuid "github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (s *Server) InitializeImagesRoutes() {
	s.Router.HandleFunc("/fill/{width}/{height}/{file_url:.*}", s.processImage).Methods("GET")
}

func (s *Server) processImage(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	url := params["file_url"]

	width, err := strconv.Atoi(params["width"])
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	height, err := strconv.Atoi(params["height"])
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	var data []byte

	key := fmt.Sprintf("%v_%v_%s", width, height, url)

	value, found := s.cache.Get(lrucache.Key(key))

	if found {
		loadedData, err := s.storage.Load(value.(string))
		if err != nil {
			respondWithMessage(writer, http.StatusInternalServerError, err.Error())
			return
		}
		data = loadedData
		s.logger.Debug(fmt.Sprintf("%s, loaded from cache", url))

		writer.Header().Set("cached", "true")
	}

	if !found {
		// load file from remote server
		loadedData, err := s.getFileFromRemoteServer(request.Context(), request.Header, url)
		if err != nil {
			writer.WriteHeader(http.StatusBadGateway)
			return
		}

		// process image
		data, err = s.saveImage(key, url, width, height, loadedData)
		if err != nil {
			statusCode := http.StatusInternalServerError
			if errors.Is(err, imageprocessor.ErrImageSmallerThanPreview) {
				statusCode = http.StatusBadRequest
			}
			respondWithMessage(writer, statusCode, err.Error())
			return
		}
	}

	if _, err := writer.Write(data); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.Header().Add("Content-Type", "image/jpeg")
	writer.Header().Set("Content-Length", strconv.Itoa(len(data)))
}

func (s *Server) getFileFromRemoteServer(ctx context.Context, header http.Header, url string) ([]byte, error) {
	fileRequest, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://"+url, nil)
	if err != nil {
		return nil, err
	}

	fileRequest.Header = header

	response, err := s.client.Do(fileRequest)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)

	return data, err
}

func (s *Server) saveImage(key, url string, width, height int, content []byte) ([]byte, error) {
	// process image
	processor := imageprocessor.NewImageProcessor()
	cropData, err := processor.Crop(width, height, content)
	if err != nil {
		return nil, err
	}

	// save image into the storage and cache
	fileName := fmt.Sprintf("%s_%v_%v_%s", guuid.New(), width, height, path.Base(url))
	if err := s.storage.Save(fileName, cropData); err != nil {
		return nil, err
	}
	s.cache.Set(lrucache.Key(key), fileName)
	s.logger.Debug(fmt.Sprintf("%s saved into the cache and storage", url))

	return cropData, nil
}

func respondWithMessage(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(code)
	w.Write([]byte(message))
}
