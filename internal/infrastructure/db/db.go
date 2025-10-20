package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"getswing.app/player-service/internal/infrastructure/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectPostgres(cfg config.Config) (*gorm.DB, *sql.DB, error) {

	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%+v/%s?sslmode=%s", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.DBSslMode)

	gormDB, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{})
	if err != nil {
		return nil, nil, fmt.Errorf("gorm open: %w", err)
	}
	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, nil, fmt.Errorf("get sql db: %w", err)
	}
	// Connection pool
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
	return gormDB, sqlDB, nil
}

func Ping(ctx context.Context, sqlDB *sql.DB) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	return sqlDB.PingContext(ctx)
}
