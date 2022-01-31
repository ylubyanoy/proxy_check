package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

// PGStore ...
type PGStore struct {
	conn            *pgxpool.Pool
	proxyRepository *ProxyRepository
}

// New is create new connect to DB
func New(dbURL string) (*PGStore, error) {
	conn, err := pgxpool.Connect(context.Background(), dbURL)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	s := &PGStore{
		conn: conn,
	}

	return s, nil
}

// Close is closing connections to DB
func (s *PGStore) Close() {
	if s.conn != nil {
		s.conn.Close()
	}
}

// Proxy ...
func (s *PGStore) Proxy() *ProxyRepository {
	if s.proxyRepository != nil {
		return s.proxyRepository
	}

	s.proxyRepository = &ProxyRepository{
		pgstore: s,
	}

	return s.proxyRepository
}
