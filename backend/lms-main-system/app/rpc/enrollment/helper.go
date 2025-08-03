package enrollment

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

func isPgError(err error, code string) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == code
	}
	return false
}
