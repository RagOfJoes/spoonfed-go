package services

import (
	"context"
	"errors"

	"github.com/RagOfJoes/spoonfed-go/internal/graphql/model"
	"github.com/RagOfJoes/spoonfed-go/internal/models"
	"github.com/RagOfJoes/spoonfed-go/pkg/auth"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

var (
	errInvalidRecipeID   = errors.New("invalid recipe id provided")
	errInvalidRecipeSlug = errors.New("invalid recipe slug provided")
)

// 								Queries
// ###########################################

// GetRecipeDetail gets a recipe's detail
func GetRecipeDetail(ctx context.Context, db *gorm.DB, slug string) (*model.Recipe, error) {
	if slug == "" {
		return nil, errInvalidRecipeSlug
	}
	tx := db.Begin()
	r := models.Recipe{}
	tx = tx.Find(&r, "slug = ?", slug)
	if err := tx.Error; err != nil {
		return nil, err
	}
	recipe, err := model.BuildRecipe(&r)
	if err != nil {
		return nil, err
	}
	return recipe, nil
}

// 								Mutations
// ###########################################

func ToggleRecipeLike(ctx context.Context, db *gorm.DB, recipeID string) (*model.Recipe, error) {
	user, err := IsAuthenticated(ctx)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, auth.ErrUnauthenticated
	}
	id, err := uuid.FromString(recipeID)
	if err != nil {
		return nil, err
	}

	r := models.Recipe{}
	find := db.Find(&r, "id = ?", id)
	if err := find.Error; err != nil {
		return nil, err
	}
	if found := find.RowsAffected; found == 0 {
		return nil, errInvalidRecipeID
	}

	recipe, err := model.BuildRecipe(&r)
	if err != nil {
		return nil, err
	}
	l := models.Like{
		UserID:     user.ID,
		EntityID:   id,
		EntityType: "recipes",
	}
	res := db.Model(&l).Where("entity_type = ? AND entity_id = ? AND user_id = ?", "recipes", recipe.ID, user.ID).Update("active", gorm.Expr("active * ?", -1))
	if res.RowsAffected == 0 {
		if err := db.Create(&l).Error; err != nil {
			return nil, err
		}
	}
	return recipe, nil
}

// CreateRecipe creates a recipe
func CreateRecipe(ctx context.Context, db *gorm.DB, input model.NewRecipeInput) (*model.Recipe, error) {
	user, err := IsAuthenticated(ctx)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, auth.ErrUnauthenticated
	}
	tx := db.Begin()
	newRecipe, err := model.NewRecipe(&input, user.ID)
	if err != nil {
		return nil, err
	}
	tx = tx.Create(newRecipe)
	if err := tx.Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	built, err := model.BuildRecipe(newRecipe)
	if err != nil {
		return nil, err
	}
	tx = tx.Commit()
	if err := tx.Error; err != nil {
		return nil, err
	}
	return built, nil
}

// 								Dataloaders
// ###########################################

// RecipeLikeDataloader creates a new dataloader
func RecipeLikeDataloader(ctx context.Context, db *gorm.DB, ids []string) ([]*int64, []error) {
	errs := []error{}
	// Execute query
	// then check for Error
	var dest []struct {
		Count    int64     `json:"count"`
		EntityID uuid.UUID `json:"entity_id"`
	}
	tx := db.Model(&models.Like{}).Select("entity_id, COUNT(*)").Group("entity_id").Where("active > 0 AND entity_id IN(?)", ids).Scan(&dest)
	if err := tx.Error; err != nil {
		errs = append(errs, err)
		return nil, errs
	}
	// Map likes
	likes := []*int64{}
	likesMap := map[string]*int64{}
	for _, val := range dest {
		likesMap[val.EntityID.String()] = &val.Count
	}
	// Append likes to array
	// in the order that they
	// were passed
	for _, id := range ids {
		count, ok := likesMap[id]
		if !ok {
			var def int64 = 0
			likes = append(likes, &def)
			continue
		}
		likes = append(likes, count)
	}
	return likes, nil
}
