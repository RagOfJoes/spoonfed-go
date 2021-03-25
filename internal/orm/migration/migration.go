package migration

import (
	"github.com/RagOfJoes/spoonfed-go/internal/models"
	"gorm.io/gorm"
)

func update(db *gorm.DB) error {
	err := db.AutoMigrate(
		&models.User{},
		&models.Recipe{},

		&models.Like{},
		&models.Image{},
	)
	if err != nil {
		return err
	}
	return nil
}

func Setup(db *gorm.DB) error {
	// Initialize the migration
	return update(db)
}
