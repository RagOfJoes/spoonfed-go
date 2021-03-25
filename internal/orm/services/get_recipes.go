package services

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/RagOfJoes/spoonfed-go/internal/graphql/model"
	"github.com/RagOfJoes/spoonfed-go/internal/models"
	"github.com/RagOfJoes/spoonfed-go/pkg/auth"
	"github.com/RagOfJoes/spoonfed-go/pkg/util"
	"gorm.io/gorm"
)

var (
	// define valid sort and filter keys
	recipeSortKeys   = []string{"name", "slug", "created_at"}
	recipeFilterKeys = []string{"name", "is_liked", "num_of_likes", "username", "created_at", "updated_at"}
	// transforms sort key to a valid sortable field
	recipeSortMap = map[string]string{
		"name": "slug",
	}
	// defines filter rules
	recipeFilterRules = map[string]Filter{
		"username": {
			Join: "User",
			Conditions: []model.FilterCondition{
				model.FilterConditionIs,
				model.FilterConditionIn,
				model.FilterConditionLike,
				model.FilterConditionILike,
				model.FilterConditionEquals,
			},
		},
		"is_liked": {
			Conditions: []model.FilterCondition{
				model.FilterConditionIs,
			},
			Query: func(ctx context.Context, db *gorm.DB, filter *model.FilterInput) (*gorm.DB, error) {
				user, err := auth.GetUserFromContext(ctx, util.ProjectContextKeys.User)
				if err != nil {
					return nil, err
				}
				u, ok := user.(*models.User)
				if uOk, err := u.IsValid(false); !ok || !uOk || err != nil {
					return nil, errors.New("you must be authenticated to use `is_liked` filter")
				}

				if filter.Value != true && filter.Value != "true" && filter.Value != false && filter.Value != "false" {
					return nil, fmt.Errorf("is_liked filter can only contain the value `true` or `false`")
				}
				op := " IS NOT NULL"
				query := "(SELECT id FROM likes WHERE likes.entity_id = recipes.id AND likes.user_id = ? AND likes.active > 0)"
				if filter.Value == false || filter.Value == "false" {
					op = " IS NULL"
				}
				db = db.Where(query+op, u.ID)
				return db, nil
			},
		},
		"num_of_likes": {
			Conditions: []model.FilterCondition{
				model.FilterConditionEquals,
				model.FilterConditionLessThan,
				model.FilterConditionLessThanEqual,
				model.FilterConditionGreaterThan,
				model.FilterConditionGreaterThanEqual,
				model.FilterConditionBetween,
			},
			Query: func(ctx context.Context, db *gorm.DB, filter *model.FilterInput) (*gorm.DB, error) {
				var value interface{}
				query := "(SELECT COALESCE(COUNT(likes.id) FILTER(WHERE entity_id = recipes.id AND active > 0), 0) FROM likes)" + filter.Condition.SQL()
				if filter.Condition == model.FilterConditionBetween {
					value = filter.Values
				} else {
					value = filter.Value
				}
				db = db.Where(query, value)
				return db, nil
			},
		},
	}
)

// GetRecipes creates a cursor paginated list of recipes
// and uses given filter to thin out results
func GetRecipes(ctx context.Context, db *gorm.DB, limit int, cursor *string, sort model.SortInput, filters []*model.FilterInput) (*model.RecipeConnection, error) {
	if err := checkLimit(limit); err != nil {
		return nil, err
	}
	sortClause, err := getSort(&sort)
	if err != nil {
		return nil, err
	} else if sortClause == "" {
		sortClause = "created_at"
	}

	var edges []*model.RecipeEdge
	recipes := []*models.Recipe{}

	tx := db.Begin()
	// Parse Filters
	filterRules := FilterRules{Inputs: filters, Keys: recipeFilterKeys, Rules: recipeFilterRules}
	tx, err = filterRules.Parse(ctx, tx)
	if err != nil {
		return nil, err
	}
	// Load cursor into query if needed
	if tx, err = loadCursor(tx, cursor, sort); err != nil {
		return nil, err
	}
	// Add the rest of the query
	tx = tx.Order(sortClause).Limit(limit + 1).Find(&recipes)
	// Check for errors
	if err := tx.Error; err != nil {
		return nil, err
	}
	// Build Edges
	for _, recipe := range recipes {
		built, err := model.BuildRecipe(recipe)
		if err != nil {
			return nil, err
		}
		encodedCursor := util.EncodeCursor(util.ToCamel(sort.Key, true), recipe)
		edges = append(edges, &model.RecipeEdge{
			Node:   built,
			Cursor: encodedCursor,
		})
	}
	// Handle new PageInfo
	newCursor := ""
	hasNextPage := len(edges) > limit
	if hasNextPage {
		newCursor = edges[len(edges)-1].Cursor
		edges = edges[:len(edges)-1]
	}
	newPageInfo := &model.PageInfo{
		Cursor:      newCursor,
		HasNextPage: hasNextPage,
	}
	// Build connection and return it
	connection := model.RecipeConnection{Edges: edges, PageInfo: newPageInfo}
	tx = tx.Commit()
	if err := tx.Error; err != nil {
		return nil, err
	}
	return &connection, nil
}

func checkLimit(limit int) error {
	if limit < 1 || limit > 100 {
		return errors.New("limit must be between the range of 1 - 100")
	}
	return nil
}

func getSort(sort *model.SortInput) (string, error) {
	flag := false
	sort.Key = util.ToSnake(sort.Key, false)
	for _, key := range recipeSortKeys {
		if key == sort.Key {
			flag = true
			break
		}
	}
	if !flag {
		return "", fmt.Errorf("invalid sort key provided. valid recipe sort keys are: %s", strings.Join(recipeSortKeys, ", "))
	}
	actualKey, ok := recipeSortMap[sort.Key]
	if ok {
		sort.Key = actualKey
	}
	sortBy := sort.Key + " " + sort.Order.String()
	return sortBy, nil
}

func loadCursor(db *gorm.DB, cursor *string, sort model.SortInput) (*gorm.DB, error) {
	decoded := ""
	if cursor != nil && len(*cursor) > 0 {
		d, err := util.DecodeCursor(*cursor)
		if err != nil {
			return nil, err
		}
		decoded = *d
	}
	if len(decoded) > 0 {
		op := " <= ?"
		if sort.Order == model.SortOrderAsc {
			op = " >= ?"
		}
		db = db.Where(sort.Key+op, decoded)
		if err := db.Error; err != nil {
			return nil, err
		}
	}
	return db, nil
}
