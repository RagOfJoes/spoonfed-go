package models

import (
	"reflect"

	"github.com/RagOfJoes/spoonfed-go/pkg/util"
	"github.com/gofrs/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Recipe struct {
	Base
	ImportURL    *string
	Name         string         `gorm:"not null;index"`
	Servings     string         `gorm:"not null;size:64"`
	Slug         string         `gorm:"not null;uniqueIndex"`
	Ingredients  pq.StringArray `gorm:"type:text[];not null"`
	Instructions pq.StringArray `gorm:"type:text[];not null"`
	Time         RecipeTime     `gorm:"embedded;embeddedPrefix:time_"`
	CreatedBy    uuid.UUID      `gorm:"not null;index"`
	// Associations
	Likes  []Like  `gorm:"polymorphic:Entity"`
	Images []Image `gorm:"polymorphic:Entity"`
	User   User    `gorm:"<-:false;foreignKey:CreatedBy"`
}

// 								Helpers
// ###########################################

// IsValid validates a Recipe
func (o *Recipe) IsValid(new bool) (ok bool, err error) {
	// Check name
	name := o.Name
	if !util.CheckLen(name, 4, 64) {
		return false, ErrInvalidLength("recipe.name", 4, 64)
	} else if !util.IsAlphaNumeric(name, true) {
		return false, ErrInvalidAlphanumeric("recipe.name")
	}
	// Check ImportURL
	url := o.ImportURL
	if url != nil && (!util.CheckLen(*url, 1, 1024) || !util.IsURL(*url)) {
		return false, ErrInvalidImportURL
	}
	// Check Servings
	if !util.CheckLen(o.Servings, 1, 64) {
		return false, ErrInvalidLength("recipe.servings", 1, 64)
	}
	// Check Ingredients and Instructions
	if l := len(o.Ingredients); l < 1 || l > 64 {
		return false, ErrInvalidElements("recipe.ingredients", 1, 64)
	}
	if l := len(o.Instructions); l < 1 || l > 64 {
		return false, ErrInvalidElements("recipe.instructions", 1, 64)
	}
	// Check Time
	t := o.Time
	if ok, err := t.IsValid(); !ok || err != nil {
		return false, err
	}
	if !new {
		if len(o.Slug) <= 0 {
			return false, ErrInvalidLength("recipe.slug", 24, 84)
		}
		if o.CreatedAt.IsZero() {
			return false, ErrRecipeInvalidCreatedAt
		}
	}
	return true, nil
}

// 									Hooks
// ###########################################

func (o *Recipe) BeforeCreate(tx *gorm.DB) (err error) {
	// validate
	if ok, err := o.IsValid(true); !ok || err != nil {
		return err
	}
	// assign a slug
	o.Slug = util.Slug(o.Name)
	return nil
}

// 							RecipeTime
// ###########################################

type RecipeTime struct {
	Prep     *string `gorm:"size:32" json:"prep"`
	Cook     *string `gorm:"size:32" json:"cook"`
	Ready    *string `gorm:"size:32" json:"ready"`
	Active   *string `gorm:"size:32" json:"active"`
	Inactive *string `gorm:"size:32" json:"inactive"`
	Total    string  `gorm:"size:32;not null" json:"total"`
}

// IsValid validates a Recipe's Time type
func (t *RecipeTime) IsValid() (ok bool, err error) {
	val := reflect.ValueOf(t).Elem()
	for i := 0; i < val.NumField(); i++ {
		name := val.Type().Field(i).Name
		err := ErrInvalidLength("recipe.time."+name, 1, 32)
		s, sOk := val.FieldByName(name).Interface().(string)
		p, pOk := val.FieldByName(name).Interface().(*string)
		if sOk {
			if !util.CheckLen(s, 1, 32) {
				return false, err
			}
		} else if pOk && p != nil {
			if !util.CheckLen(*p, 1, 32) {
				return false, err
			}
		}
	}
	return true, nil
}
