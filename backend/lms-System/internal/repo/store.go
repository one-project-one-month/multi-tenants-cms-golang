package repo

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store interface {
	Querier
	RegisterUserWithRoles(ctx context.Context, arg RegisterUserWithRolesParams) (RegisterUserWithRolesRow, error)
	GetUserByEmail(ctx context.Context, lmsUserEmail string) (LmsUser, error)
	UpdateEmailVerification(ctx context.Context, lmsUserEmail string) error
}
type SQLStore struct {
	*Queries
	pool *pgxpool.Pool
}

func NewStore(pool *pgxpool.Pool) Store {
	return &SQLStore{
		Queries: New(pool),
		pool:    pool,
	}
}
