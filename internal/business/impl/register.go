package impl

import (
	"context"
	"fmt"

	"github.com/chernyshevuser/gopfermart.git/internal/business/impl/utils"
	"github.com/chernyshevuser/gopfermart.git/internal/db"
	query "github.com/chernyshevuser/gopfermart.git/internal/db/impl"
	"github.com/jackc/pgx/v5"
)

func (g *gophermart) Register(ctx context.Context, login, password string) (ok bool, err error) {
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
		return false, fmt.Errorf("creating db tx failed: %w", err)
	}
	defer tx.Rollback(ctx)

	err = db.SimpleInTx(ctx, tx, func(ctx context.Context, tx pgx.Tx) (err error) {
		exists, err := query.IsUserIxist(ctx, tx, login)
		if err != nil {
			return err
		}

		ok = true
		if exists {
			ok = false

			p, err := query.GetUserPassword(ctx, tx, login)
			if err != nil {
				return err
			}

			g.storage.Set(login, p)

			return nil
		}

		return query.Register(ctx, tx, login, encryptedPassword)
	})
	if err != nil {
		return ok, err
	}

	if err = tx.Commit(ctx); err != nil {
		return ok, fmt.Errorf("error in commiting transaction: %w", err)
	}

	return ok, nil
}
