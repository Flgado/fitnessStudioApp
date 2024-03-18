package database

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

func Transaction(ctx context.Context, conn sqlx.DB, f func(tx *sqlx.Tx) error) error {
	tx, err := conn.BeginTxx(ctx, nil)

	if err != nil {
		return fmt.Errorf("Begin %w", err)
	}

	if err := f(tx); err != nil {
		_ = tx.Rollback()

		return fmt.Errorf("f %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("Commit %w", err)
	}

	return nil
}
