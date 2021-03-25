package model

import (
	"errors"

	"github.com/RagOfJoes/spoonfed-go/internal/models"
	"github.com/RagOfJoes/spoonfed-go/pkg/util"
	"github.com/gofrs/uuid"
)

var (
	ErrFailedToBuildRecipe  = errors.New("failed to build recipe")
	ErrFailedToEscapeRecipe = errors.New("failed to escape recipe")
)

type Recipe struct {
	ID           string             `json:"id"`
	Name         string             `json:"name"`
	Slug         string             `json:"slug"`
	Servings     string             `json:"servings"`
	ImportURL    *string            `json:"importUrl"`
	CreatedBy    string             `json:"createdBy"`
	Ingredients  []string           `json:"ingredients"`
	Instructions []string           `json:"instructions"`
	Date         *MetaDate          `json:"date"`
	Time         *models.RecipeTime `json:"time"`
}

type NewRecipeInput struct {
	Name         string           `json:"name"`
	Servings     string           `json:"servings"`
	ImportURL    *string          `json:"importUrl"`
	Time         *RecipeTimeInput `json:"time"`
	Images       []*ImageInput    `json:"images"`
	Ingredients  []string         `json:"ingredients"`
	Instructions []string         `json:"instructions"`
}

func BuildRecipe(r *models.Recipe) (*Recipe, error) {
	ok, err := r.IsValid(false)
	if !ok || err != nil {
		return nil, err
	}
	util.UnescapeStruct(r)
	recipe := &Recipe{
		ID:           r.ID.String(),
		Name:         r.Name,
		Slug:         r.Slug,
		Servings:     r.Servings,
		ImportURL:    r.ImportURL,
		CreatedBy:    r.CreatedBy.String(),
		Ingredients:  r.Ingredients,
		Instructions: r.Instructions,
		Time:         &r.Time,
		Date: &MetaDate{
			Creation:   r.CreatedAt,
			LastUpdate: r.UpdatedAt,
		},
	}
	return recipe, nil
}

func NewRecipe(r *NewRecipeInput, createdBy uuid.UUID) (*models.Recipe, error) {
	images := []models.Image{}
	for _, image := range r.Images {
		n, err := NewImage(image)
		if err != nil {
			return nil, err
		}
		images = append(images, *n)
	}
	util.EscapeStruct(r)
	newRecipe := &models.Recipe{
		Name:         r.Name,
		Servings:     r.Servings,
		ImportURL:    r.ImportURL,
		Ingredients:  r.Ingredients,
		Instructions: r.Instructions,
		CreatedBy:    createdBy,
		Time: models.RecipeTime{
			Prep:     r.Time.Prep,
			Cook:     r.Time.Cook,
			Ready:    r.Time.Ready,
			Active:   r.Time.Active,
			Inactive: r.Time.Inactive,
			Total:    r.Time.Total,
		},
		Images: images,
	}
	isOk, err := newRecipe.IsValid(true)
	if !isOk || err != nil {
		return nil, err
	}
	return newRecipe, nil
}
