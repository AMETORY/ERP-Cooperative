package services

import (
	"ametory-cooperative/config"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB(cfg *config.Config) (*gorm.DB, error) {
	var dialector gorm.Dialector

	switch cfg.Database.Type {
	case "postgres":
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
			cfg.Database.Host, cfg.Database.User, cfg.Database.Password, cfg.Database.Name, cfg.Database.Port)
		dialector = postgres.Open(dsn)
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.Name)
		dialector = mysql.Open(dsn)
	case "sqlite":
		dialector = sqlite.Open(cfg.Database.Name) // Name adalah path ke file SQLite
	case "mssql":
		dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s",
			cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.Name)
		dialector = sqlserver.Open(dsn)
	default:
		return nil, fmt.Errorf("database type tidak didukung: %s", cfg.Database.Type)
	}
	options := []gorm.Option{}
	if config.App.Server.Debug {
		options = append(options, &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	}

	db, err := gorm.Open(dialector, options...)

	if err != nil {
		return nil, fmt.Errorf("gagal menghubungkan ke database: %v", err)
	}

	return db, nil
}
