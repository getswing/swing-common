package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	common_logger "getswing.app/player-service/pkg/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type CommonDBConfig struct {
	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     string
	DBName     string
	DBSslMode  string
}

func ConnectPostgresDB(cfg CommonDBConfig, serviceName string) (*gorm.DB, *sql.DB, error) {
	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%+v/%s?sslmode=%s", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.DBSslMode)

	gormDB, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{
		Logger: common_logger.NewGormLogger(serviceName),
	})
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
