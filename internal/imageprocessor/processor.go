package imageprocessor

import (
	"bytes"
	"errors"
	"image/jpeg"

	"github.com/disintegration/imaging"
)

var ErrImageSmallerThanPreview = errors.New("source image smaller than preview")

type ImageProcessor interface {
	Crop(width, height int, data []byte) ([]byte, error)
}

type processor struct{}

func NewImageProcessor() ImageProcessor {
	return &processor{}
}

func (p *processor) Crop(width, height int, data []byte) ([]byte, error) {
	sourceImage, err := imaging.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	if sourceImage.Bounds().Dx() < width || sourceImage.Bounds().Dy() < height {
		return nil, ErrImageSmallerThanPreview
	}

	sourceImage = imaging.Fill(sourceImage, width, height, imaging.Center, imaging.Lanczos)

	var buff bytes.Buffer

	err = jpeg.Encode(&buff, sourceImage, nil)

	return buff.Bytes(), err
}
