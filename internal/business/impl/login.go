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

func (g *gophermart) Login(ctx context.Context, login, password string) (ok bool, sessionToken string, err error) {
	ok = true

	encryptedPassword, exists := g.storage.Get(login)
	if !exists {
		tx, err := g.db.BeginR(ctx)
		if err != nil {
			return false, "", fmt.Errorf("creating db tx failed: %w", err)
		}
		defer tx.Rollback(ctx)

		err = db.SimpleInTx(ctx, tx, func(ctx context.Context, tx pgx.Tx) (err error) {
			encryptedPassword, err = query.GetUserPassword(ctx, tx, login)
			if err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					ok = false
					return nil
				}
				return err
			}

			g.storage.Set(login, encryptedPassword)
			return nil
		})
		if err != nil {
			return false, "", err
		}

		if err = tx.Commit(ctx); err != nil {
			return false, "", fmt.Errorf("error in commiting transaction: %w", err)
		}
	}

	if !ok {
		return false, "", nil
	}

	decryptedPassword, err := utils.Decrypt(encryptedPassword)
	if err != nil {
		g.logger.Errorw(
			"decryption failed",
			"msg", err,
		)
		return ok, "", err
	}

	ok = decryptedPassword == password
	if !ok {
		return ok, "", nil
	}

	sessionToken, err = g.sessionSvc.NewToken(login)
	if err != nil {
		return ok, "", err
	}

	return ok, sessionToken, nil
}
