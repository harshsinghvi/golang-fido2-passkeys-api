package helpers

import (
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
)

func PgErrorCodeAndMessage(err error) (string, string) {
	// log.Print(err.(*pgconn.PgError).Code)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code, pgErr.Message
	}
	return "", ""
}
