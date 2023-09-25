package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/ilya-rusyanov/shrinklator/internal/entities"
	"github.com/ilya-rusyanov/shrinklator/internal/logger"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
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

func (p *Postgres) MustClose() {
	err := p.db.Close()

	if err != nil {
		panic(fmt.Errorf("error closing DB: %w", err))
	}
}

func (p *Postgres) Put(ctx context.Context, id, value string, uid *entities.UserID) error {
	_, err := p.db.ExecContext(ctx,
		`INSERT INTO shorts (short, long, user_id) VALUES ($1, $2, $3)`, id, value, uid)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			row := p.db.QueryRowContext(ctx,
				`SELECT short FROM shorts WHERE long = $1`, value)
			var short string
			row.Scan(&short)
			if err := row.Err(); err != nil {
				return fmt.Errorf("error scanning value: %w", err)
			}

			return ErrAlreadyExists{
				StoredValue: short,
			}
		}
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

	stmt, err := tx.PrepareContext(ctx, `INSERT INTO shorts (short, long)
VALUES ($1, $2)`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, pair := range data {
		_, err := stmt.ExecContext(ctx, pair.Short, pair.Long)
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

func (p *Postgres) ByID(ctx context.Context, id string) (entities.ExpandResult, error) {
	row := p.db.QueryRowContext(ctx,
		`SELECT long, is_deleted FROM shorts WHERE short = $1`, id)
	var res entities.ExpandResult
	row.Scan(&res.URL, &res.Removed)
	if err := row.Err(); err != nil {
		return res, fmt.Errorf("error fetching record: %w", err)
	}
	p.log.Debug("successfull record fetch")
	return res, nil
}

func (p *Postgres) ByUID(ctx context.Context,
	uid entities.UserID) (entities.PairArray, error) {
	p.log.Info("selecting by uid", zap.String("uid", string(uid)))
	rows, err := p.db.QueryContext(ctx,
		`SELECT short, long FROM shorts WHERE user_id = $1`, uid)
	if err != nil {
		return nil, fmt.Errorf("error selecting rows: %w", err)
	}
	defer rows.Close()

	var pairs entities.PairArray

	for rows.Next() {
		var pair entities.ShortLongPair
		err = rows.Scan(&pair.Short, &pair.Long)

		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}

		pairs = append(pairs, pair)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	p.log.Info("success", zap.Int("rows found", len(pairs)))

	return pairs, nil
}

func (p *Postgres) Delete(ctx context.Context, req entities.DeleteRequest) error {
	values := make([]string, len(req))
	for i, r := range req {
		values[i] = fmt.Sprintf("( short = '%s' AND user_id = '%s' )", r.URL, r.UID)
	}
	cmd :=
		"UPDATE shorts SET is_deleted = true WHERE " + strings.Join(values, " OR ") + ";"
	err, _ := p.db.ExecContext(ctx, cmd)
	if err != nil {
		return fmt.Errorf("error marking for deletion: %w", err)
	}
	return nil
}

func migrate(ctx context.Context, log *logger.Log, db *sql.DB) error {
	_, err := db.ExecContext(ctx,
		`CREATE TABLE IF NOT EXISTS shorts
(short text, long text UNIQUE, user_id text, is_deleted boolean, PRIMARY KEY (short))`)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}
	log.Info("db migrated")

	return nil
}
