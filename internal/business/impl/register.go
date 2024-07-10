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

func (g *gophermart) Register(ctx context.Context, login, password string) (ok bool, sessionToken string, err error) {
	ok = true

	_, exists := g.storage.Get(login)
	if exists {
		ok = false
		return
	}

	encryptedPassword, err := utils.Encrypt(password)
	if err != nil {
		g.logger.Errorw(
			"decryption failed",
			"msg", err,
		)
		return
	}

	tx, err := g.db.BeginW(ctx)
	if err != nil {
		return false, "", fmt.Errorf("creating db tx failed: %w", err)
	}
	defer tx.Rollback(ctx)

	err = db.SimpleInTx(ctx, tx, func(ctx context.Context, tx pgx.Tx) (err error) {
		_, err = query.GetUserPassword(ctx, tx, login)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return query.Register(ctx, tx, login, encryptedPassword)
			}
			return err
		}

		ok = false
		return nil
	})
	if err != nil {
		return false, "", err
	}

	if err = tx.Commit(ctx); err != nil {
		return ok, "", fmt.Errorf("error in commiting transaction: %w", err)
	}

	if !ok {
		return false, "", nil
	}

	sessionToken, err = g.sessionSvc.NewToken(login)
	if err != nil {
		return false, "", err
	}

	return true, sessionToken, nil
}
