package db

import (
	"database/sql" // Required for sql.ErrNoRows
	"errors"       // Required for errors.Is
	"fmt"
	// "strings" // No longer needed for this logic
)

// SaveConfig saves or updates the single application configuration row (id=1).
// This function uses an "UPSERT" pattern.
func (db *DB) SaveConfig(cfg Config) error {
	// The saveConfig query expects 11 arguments:
	// 1 (for id=1), then the 10 fields from the Config struct.
	_, err := db.SQL.ExecContext(db.ctx, saveConfig,
		1, // id (hardcoded to 1 for the UPSERT)
		cfg.Port,
		cfg.EnableRateLimit,
		cfg.RateLimit,
		cfg.OutputFolder,
		cfg.LoopInterval,
		cfg.EnableGDrive,
		cfg.GDriveFilepath,
		cfg.TelegramChatID,
		cfg.TelegramToken,
		cfg.EnableTelegram,
	)

	if err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	return nil
}

// GetConfig retrieves the single application configuration row (id=1).
func (db *DB) Config() (Config, error) {
	var cfg Config

	// Use QueryRowContext since we only ever expect one row.
	// The getConfig query takes no parameters.
	row := db.SQL.QueryRowContext(db.ctx, getConfig)

	err := row.Scan(
		&cfg.Port,
		&cfg.EnableRateLimit,
		&cfg.RateLimit,
		&cfg.OutputFolder,
		&cfg.LoopInterval,
		&cfg.EnableGDrive,
		&cfg.GDriveFilepath,
		&cfg.TelegramChatID,
		&cfg.TelegramToken,
		&cfg.EnableTelegram,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// If no row exists yet, return a default/empty Config struct.
			// The application can then use this to insert the first config.
			return Config{}, nil
		}
		// A real database error occurred.
		return Config{}, fmt.Errorf("failed to scan config: %w", err)
	}

	return cfg, nil
}
