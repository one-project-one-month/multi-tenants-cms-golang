package db

import "github.com/jackc/pgx/v5/pgxpool"

type Store struct {
	*Queries
	conn *pgxpool.Pool
}

func NewStore(conn *pgxpool.Pool) *Store {
	return &Store{
		conn: conn,
	}
}
