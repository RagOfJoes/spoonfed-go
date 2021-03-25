package model

import (
	"errors"

	"github.com/RagOfJoes/spoonfed-go/internal/models"
)

var (
	ErrFailedToEscapeImage = errors.New("failed escape image input")
)

type Image struct {
	URL     string `json:"url"`
	Name    string `json:"name"`
	Caption string `json:"caption"`
}

type ImageInput struct {
	URL     string `json:"url"`
	Name    string `json:"name"`
	Caption string `json:"caption"`
}

func BuildImage(i *models.Image) (*Image, error) {
	ok, err := i.IsValid(false)
	if !ok || err != nil {
		return nil, err
	}
	image := &Image{
		URL:     i.URL,
		Name:    i.Name,
		Caption: i.Caption,
	}
	return image, nil
}

func NewImage(i *ImageInput) (*models.Image, error) {
	return &models.Image{
		URL:     i.URL,
		Name:    i.Name,
		Caption: i.Caption,
	}, nil
}
