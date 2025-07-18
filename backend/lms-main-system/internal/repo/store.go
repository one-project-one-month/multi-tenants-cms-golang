package repo

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	authenticationpb "github.com/multi-tenants-cms-golang/lms-sys/protogen/authentication"
	"github.com/sirupsen/logrus"
)

type Store interface {
	Querier
	ExecTx(
		ctx context.Context,
		fn func(*Queries) error,
	) error
	WithTx(tx pgx.Tx) *Queries
	UpdateEmailVerification(
		ctx context.Context,
		lmsUserEmail string,
	) error
	WholeRegistrationFlow(
		ctx context.Context,
		req *authenticationpb.RegisterRequest,
		namespace string,
	) (*authenticationpb.RegisterResponse, error)
}
type SQLStore struct {
	*Queries
	logger   *logrus.Logger
	connPool *pgxpool.Pool
}

func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.connPool.Begin(ctx)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit(ctx)
}

func (store *SQLStore) ExecTx(ctx context.Context, fn func(*Queries) error) error {
	return store.execTx(ctx, fn)
}

func NewStore(
	logger *logrus.Logger,
	pool *pgxpool.Pool,
) Store {
	return &SQLStore{
		logger:   logger,
		Queries:  New(pool),
		connPool: pool,
	}
}
