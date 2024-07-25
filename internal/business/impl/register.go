package impl

import (
	"context"
	"errors"
	"fmt"

	"github.com/chernyshevuser/gopfermart.git/internal/business/impl/utils"
	"github.com/chernyshevuser/gopfermart.git/internal/db"
	query "github.com/chernyshevuser/gopfermart.git/internal/db/impl"
	"github.com/jackc/pgx/v5"
)

// Register returns nil sessionToken if user is already exist
func (g *gophermart) Register(ctx context.Context, login, password string) (sessionToken *string, err error) {
	encryptedPassword, err := utils.Encrypt(password)
	if err != nil {
		g.logger.Errorw(
			"pass encryption failed",
			"login", login,
			"msg", err,
		)
		return
	}

	tx, err := g.db.BeginW(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create db tx: %w", err)
	}
	defer tx.Rollback(ctx)

	err = db.SimpleInTx(ctx, tx, func(ctx context.Context, tx pgx.Tx) (err error) {
		_, err = query.GetEncryptedPassword(ctx, tx, login)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return query.Register(ctx, tx, login, encryptedPassword)
			}

			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit db tx: %w", err)
	}

	token, err := g.sessionSvc.NewToken(login)
	if err != nil {
		return nil, fmt.Errorf("failed to create new session token: %w", err)
	}

	return &token, nil
}
