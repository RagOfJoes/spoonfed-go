package models

import (
	"errors"
	"fmt"
)

var (
	// Common
	ErrInvalidAlphanumeric = func(key string) error {
		return fmt.Errorf("%s must not contain any special characters", key)
	}
	ErrInvalidLength = func(key string, min int, max int) error {
		return fmt.Errorf("%s must be between %d and %d characters long", key, min, max)
	}
	ErrInvalidElements = func(key string, min int, max int) error {
		return fmt.Errorf("%s must contain between %d and %d elements", key, min, max)
	}

	// Image
	ErrImageInvalidURL = errors.New("image url is invalid. make sure url begins with http/https")

	// User
	ErrUserNotBuilt  = errors.New("proper fields have not been set")
	ErrInvalidEmail  = errors.New("invalid email provided")
	ErrInvalidAvatar = errors.New("invalid user avatar url provided. make sure url begins with http/https")

	// Recipes
	ErrRecipeNotBuilt         = errors.New("proper fields have not been set")
	ErrRecipeInvalidCreatedAt = errors.New("invalid created at time provided")
	ErrInvalidTime            = errors.New("recipe must contain a valid time type")
	ErrInvalidImportURL       = errors.New("recipe import url is invalid. make sure url begins with http/https")
)
