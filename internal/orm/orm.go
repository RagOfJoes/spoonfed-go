package orm

import (
	"github.com/RagOfJoes/spoonfed-go/internal/orm/migration"
	"github.com/RagOfJoes/spoonfed-go/pkg/logger"
	"github.com/RagOfJoes/spoonfed-go/pkg/util"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ORM struct {
	DB *gorm.DB
}

func New(cfg *util.ServerConfig) (*ORM, error) {
	db, err := gorm.Open(postgres.Open(cfg.ORM.DSN), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	orm := &ORM{DB: db}
	if cfg.ORM.AutoMigrate {
		err = migration.Setup(orm.DB)
		if err != nil {
			return nil, err
		}
	}
	logger.Info("[DB] Connection Initialized.")
	return orm, nil
}
