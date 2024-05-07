package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

// VerifyEmailTxParams contains the input parameters of verifying email transaction
type VerifyEmailTxParams struct {
	EmailId int64
	Token   string
}

// VerifyEmailTxResult is the result of verifying email transaction
type VerifyEmailTxResult struct {
	User        User
	VerifyEmail VerifyEmail
}

// VerifyEmailTx verify email within transaction
func (store *SQLStore) VerifyEmailTx(ctx context.Context, arg VerifyEmailTxParams) (VerifyEmailTxResult, error) {
	var result VerifyEmailTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		verifyEmail, err := q.UpdateVerifyEmail(ctx, UpdateVerifyEmailParams{
			ID:    arg.EmailId,
			Token: arg.Token,
		})
		result.VerifyEmail = verifyEmail
		if err != nil {
			return err
		}
		user, err := q.UpdateUser(ctx, UpdateUserParams{
			Username: verifyEmail.Username,
			IsEmailVerified: pgtype.Bool{
				Bool:  true,
				Valid: true,
			},
		})
		if err != nil {
			return err
		}
		result.User = user

		return err
	})

	return result, err
}
