package model


type Recipe struct {
	Name         string             `json:"name"`
	Servings     string             `json:"servings"`
	Time         *RecipeTimeType    `json:"time"`
	Ingredients  []string           `json:"ingredients"`
	Instructions []string           `json:"instructions"`
	Slug         string             `json:"slug"`
	ImportURL    *string            `json:"importUrl"`
	Images       []*Image           `json:"images"`
	Date         *MetaDate          `json:"date"`
}
