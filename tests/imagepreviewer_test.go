//go:build integration
// +build integration

package tests

import (
	"bytes"
	"context"
	"fmt"
	"image"
	_ "image/jpeg"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	client *http.Client
}

func NewTestSuite() *TestSuite {
	return &TestSuite{client: http.DefaultClient}
}

func (s *TestSuite) DoGetRequest(t *testing.T, width, height int, url string) (*http.Response, []byte, error) {
	t.Helper()
	requestURL := fmt.Sprintf("http://image-previewer:80/fill/%v/%v/%s", width, height, url)

	req, err := http.NewRequestWithContext(context.Background(), "GET", requestURL, nil)
	require.NoError(t, err)

	response, err := s.client.Do(req)
	require.NoError(t, err)

	defer response.Body.Close()

	b, err := io.ReadAll(response.Body)

	require.NoError(t, err)

	return response, b, err
}

func TestImagePreviewer(t *testing.T) {
	s := NewTestSuite()

	//nolint:bodyclose
	_, _, err := s.DoGetRequest(t, 500, 400, "nginx/_gopher_original_1024x504.jpg")

	require.NoError(t, err)
}

func TestServerDoesntExist(t *testing.T) {
	s := NewTestSuite()

	//nolint:bodyclose
	response, _, err := s.DoGetRequest(t, 2000, 400, "wrong-server.com/_gopher_original_1024x504.jpg")

	require.NoError(t, err)

	require.Equal(t, http.StatusBadGateway, response.StatusCode)
}

func TestWrongImageSize(t *testing.T) {
	s := NewTestSuite()

	//nolint:bodyclose
	response, body, err := s.DoGetRequest(t, 2000, 400, "nginx/_gopher_original_1024x504.jpg")

	require.NoError(t, err)
	require.Equal(t, "text/plain", response.Header.Get("Content-Type"))
	require.Equal(t, "source image smaller than preview", string(body))
}

func TestWrongFileFormat(t *testing.T) {
	s := NewTestSuite()

	//nolint:bodyclose
	response, _, err := s.DoGetRequest(t, 500, 400, "nginx/test.txt")

	require.NoError(t, err)
	require.Equal(t, http.StatusInternalServerError, response.StatusCode)
}

func TestSuccessCrop(t *testing.T) {
	s := NewTestSuite()

	//nolint:bodyclose
	response, body, err := s.DoGetRequest(t, 500, 400, "nginx/_gopher_original_1024x504.jpg")

	require.NoError(t, err)
	require.Equal(t, response.StatusCode, http.StatusOK)
	require.Equal(t, "image/jpeg", response.Header.Get("Content-Type"))

	println(body)
	resultImage, _, err := image.Decode(bytes.NewReader(body))

	require.NoError(t, err)

	require.Equal(t, resultImage.Bounds().Dx(), 500)
	require.Equal(t, resultImage.Bounds().Dy(), 400)
}

func TestGetImageFromCache(t *testing.T) {
	s := NewTestSuite()

	//nolint:bodyclose
	response, body, err := s.DoGetRequest(t, 500, 450, "nginx/_gopher_original_1024x504.jpg")

	require.NoError(t, err)
	require.Equal(t, response.StatusCode, http.StatusOK)
	require.Equal(t, "image/jpeg", response.Header.Get("Content-Type"))
	require.Equal(t, "", response.Header.Get("cached"))

	resultImage, _, err := image.Decode(bytes.NewReader(body))
	require.NoError(t, err)
	require.Equal(t, resultImage.Bounds().Dx(), 500)
	require.Equal(t, resultImage.Bounds().Dy(), 450)

	//nolint:bodyclose
	response1, body1, err1 := s.DoGetRequest(t, 500, 450, "nginx/_gopher_original_1024x504.jpg")

	require.NoError(t, err1)
	require.Equal(t, response1.StatusCode, http.StatusOK)
	require.Equal(t, "image/jpeg", response1.Header.Get("Content-Type"))
	require.Equal(t, "true", response1.Header.Get("cached"))

	resultImage, _, err = image.Decode(bytes.NewReader(body1))
	require.NoError(t, err)
	require.Equal(t, resultImage.Bounds().Dx(), 500)
	require.Equal(t, resultImage.Bounds().Dy(), 450)
}
