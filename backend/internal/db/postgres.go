package db

import (
	"fmt"

	"nexus-messenger/backend/config"
	"nexus-messenger/backend/internal/channel"
	"nexus-messenger/backend/internal/message"
	"nexus-messenger/backend/internal/user"
	"nexus-messenger/backend/internal/workspace"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewPostgres(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.PostgresHost, cfg.PostgresPort,
		cfg.PostgresUser, cfg.PostgresPassword,
		cfg.PostgresDB, cfg.PostgresSSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, fmt.Errorf("postgres connect: %w", err)
	}

	if err := db.AutoMigrate(
		&user.User{},
		&workspace.Workspace{},
		&workspace.WorkspaceMember{},
		&channel.Channel{},
		&channel.ChannelMember{},
		&message.Message{},
	); err != nil {
		return nil, fmt.Errorf("postgres migrate: %w", err)
	}

	return db, nil
}
