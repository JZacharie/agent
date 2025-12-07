package database

import (
	"fmt"
	"path"

	"github.com/eduardooliveira/stLib/core/runtime"
	"github.com/glebarez/sqlite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase() error {
	var err error
	var dialector gorm.Dialector

	// Determine database type from configuration
	dbType := runtime.Cfg.Database.Type
	if dbType == "" {
		dbType = "sqlite" // Default to SQLite for backward compatibility
	}

	switch dbType {
	case "postgres":
		// Build PostgreSQL DSN
		dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			runtime.Cfg.Database.Postgres.Host,
			runtime.Cfg.Database.Postgres.Port,
			runtime.Cfg.Database.Postgres.User,
			runtime.Cfg.Database.Postgres.Password,
			runtime.Cfg.Database.Postgres.Database,
			runtime.Cfg.Database.Postgres.SSLMode,
		)
		dialector = postgres.Open(dsn)
	case "sqlite":
		// Use SQLite with file-based storage
		dialector = sqlite.Open(path.Join(runtime.GetDataPath(), "data.db"))
	default:
		return fmt.Errorf("unsupported database type: %s", dbType)
	}

	DB, err = gorm.Open(dialector, &gorm.Config{
		TranslateError: true,
	})
	if err != nil {
		return fmt.Errorf("failed to connect database: %w", err)
	}

	if err = initTags(); err != nil {
		return err
	}

	if err = initProjects(); err != nil {
		return err
	}

	if err = initAssets(); err != nil {
		return err
	}

	return nil
}
