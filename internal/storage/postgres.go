package storage

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ilya-rusyanov/shrinklator/internal/entities"
	"github.com/ilya-rusyanov/shrinklator/internal/logger"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Postgres struct {
	db  *sql.DB
	log *logger.Log
}

func NewPostgres(ctx context.Context, log *logger.Log, dsn string) (*Postgres, error) {
	db, err := sql.Open("pgx", dsn)

	if err != nil {
		return nil, fmt.Errorf("failed to open DB: %w", err)
	}
	log.Info("connected to database")

	err = migrate(ctx, log, db)
	if err != nil {
		return nil, fmt.Errorf("failed to migrate: %w", err)
	}

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

func (p *Postgres) Put(ctx context.Context, id, value string) error {
	_, err := p.db.ExecContext(ctx,
		`INSERT INTO shorts (short, long) VALUES ($1, $2)`, id, value)
	if err != nil {
		return fmt.Errorf("error writing to DB: %w", err)
	}
	p.log.Debug("successfull write to database")
	return nil
}

func (p *Postgres) PutBatch(ctx context.Context, data []entities.ShortLongPair) error {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}
	defer tx.Rollback()

	for _, pair := range data {
		_, err := tx.ExecContext(ctx, `INSERT INTO shorts (short, long)
VALUES ($1, $2)`, pair.Short, pair.Long)
		if err != nil {
			return fmt.Errorf("failed to execute statement in transaction: %w", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (p *Postgres) ByID(ctx context.Context, id string) (string, error) {
	row := p.db.QueryRowContext(ctx,
		`SELECT long FROM shorts WHERE short = $1`, id)
	var res string
	row.Scan(&res)
	fmt.Println(res)
	if err := row.Err(); err != nil {
		return "", fmt.Errorf("error fetching record: %w", err)
	}
	p.log.Debug("successfull record fetch")
	return res, nil
}

func migrate(ctx context.Context, log *logger.Log, db *sql.DB) error {
	_, err := db.ExecContext(ctx,
		`CREATE TABLE IF NOT EXISTS shorts (short text, long text UNIQUE,
PRIMARY KEY (short)
)`)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}
	log.Info("db migrated")

	return nil
}
