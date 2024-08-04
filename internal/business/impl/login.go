package impl

import (
	"context"
	"errors"
	"fmt"

	"github.com/chernyshevuser/gopfermart/internal/business"
	"github.com/chernyshevuser/gopfermart/internal/db"
	query "github.com/chernyshevuser/gopfermart/internal/db/impl"
	"github.com/jackc/pgx/v5"
)

// Login returns nil sessionToken if user/password data doesn't exist
func (g *gophermart) Login(ctx context.Context, login, password string) (sessionToken string, err error) {
	tx, err := g.db.BeginR(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create db tx: %w", err)
	}
	defer tx.Rollback(ctx)

	var encryptedPassword string

	err = db.SimpleInTx(ctx, tx, func(ctx context.Context, tx pgx.Tx) (err error) {
		encryptedPassword, err = query.GetEncryptedPassword(ctx, tx, login)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil
			}
			return err
		}

		return nil
	})
	if err != nil {
		return "", err
	}

	if err = tx.Commit(ctx); err != nil {
		return "", fmt.Errorf("failed to commit db tx: %w", err)
	}

	// wrong login/password pair
	if len(encryptedPassword) == 0 {
		return "", business.ErrUnauthorized
	}

	decryptedPassword, err := g.cryptoSvc.Decrypt(encryptedPassword)
	if err != nil {
		g.logger.Errorw(
			"failed to decrypt password",
			"msg", err,
		)
		return "", fmt.Errorf("failed to decrypt password: %w", err)
	}

	// wrong login/password pair
	if decryptedPassword != password {
		return "", business.ErrUnauthorized
	}

	token, err := g.sessionSvc.NewToken(login)
	if err != nil {
		return "", fmt.Errorf("failed to create new session token: %w", err)
	}

	return token, nil
}
