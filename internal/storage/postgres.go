package storage

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ilya-rusyanov/shrinklator/internal/logger"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Postgres struct {
	db  *sql.DB
	log *logger.Log
}

func NewPostgres(log *logger.Log, dsn string) (*Postgres, error) {
	db, err := sql.Open("pgx", dsn)

	if err != nil {
		return nil, fmt.Errorf("failed to open DB: %w", err)
	}
	log.Info("connected to database")

	return &Postgres{
		db:  db,
		log: log,
	}, nil
}

func (p *Postgres) Ping(ctx context.Context) error {
	err := p.db.PingContext(ctx)

	if err != nil {
		p.log.Warn("context ping failure")
		return fmt.Errorf("failed to ping context: %w", err)
	}

	return nil
}

func (p *Postgres) Close() error {
	err := p.db.Close()

	if err != nil {
		return fmt.Errorf("error closing DB: %w", err)
	}

	return nil
}

func (p *Postgres) Put(id, value string) error {
	return nil
}

func (p *Postgres) ByID(id string) (string, error) {
	return "", nil
}
