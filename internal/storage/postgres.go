package storage

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"strings"

	"github.com/ilya-rusyanov/shrinklator/internal/entities"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

// Postgres - postgres DB storage
type Postgres struct {
	db  *sql.DB
	log Logger
}

// NewPostgres - constructs Postgres object
func NewPostgres(ctx context.Context, log Logger, dsn string) (*Postgres, error) {
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

// Ping checks database availability
func (p *Postgres) Ping(ctx context.Context) error {
	err := p.db.PingContext(ctx)

	if err != nil {
		p.log.Warn("context ping failure")
		return fmt.Errorf("failed to ping context: %w", err)
	}

	return nil
}

// MustClose finalizes database or panics
func (p *Postgres) MustClose() {
	err := p.db.Close()

	if err != nil {
		panic(fmt.Errorf("error closing DB: %w", err))
	}
}

// Put adds entry
func (p *Postgres) Put(ctx context.Context, id, value string, uid *entities.UserID) error {
	_, err := p.db.ExecContext(ctx,
		`INSERT INTO shorts (short, long, user_id) VALUES ($1, $2, $3)`, id, value, uid)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			row := p.db.QueryRowContext(ctx,
				`SELECT short FROM shorts WHERE long = $1`, value)
			var short string
			if e := row.Scan(&short); e != nil {
				return fmt.Errorf("error scanning value: %w", e)
			}
			if e := row.Err(); e != nil {
				return fmt.Errorf("error scanning value: %w", e)
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

// PutBatch adds multiple entries
func (p *Postgres) PutBatch(ctx context.Context, data []entities.ShortLongPair) error {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}
	defer func() {
		if e := tx.Rollback(); e != nil {
			p.log.Warnf("rollback error: %q", e.Error())
		}
	}()

	stmt, err := tx.PrepareContext(ctx, `INSERT INTO shorts (short, long)
VALUES ($1, $2)`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer func() {
		if e := stmt.Close(); e != nil {
			p.log.Warnf("stmt close error: %q", e.Error())
		}
	}()

	for _, pair := range data {
		_, e := stmt.ExecContext(ctx, pair.Short, pair.Long)
		if e != nil {
			return fmt.Errorf("failed to execute statement in transaction: %w", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// ByID searches entry by identifier
func (p *Postgres) ByID(ctx context.Context, id string) (entities.ExpandResult, error) {
	var res entities.ExpandResult

	row := p.db.QueryRowContext(ctx,
		`SELECT long, is_deleted FROM shorts WHERE short = $1`, id)
	if err := row.Scan(&res.URL, &res.Removed); err != nil {
		return res, fmt.Errorf("row scan error: %w", err)
	}
	if err := row.Err(); err != nil {
		return res, fmt.Errorf("error fetching record: %w", err)
	}
	p.log.Debug("successfull record fetch")
	return res, nil
}

// ByUID searches entries by user identifier
func (p *Postgres) ByUID(ctx context.Context,
	uid entities.UserID) (entities.PairArray, error) {
	p.log.Info("selecting by uid", zap.String("uid", string(uid)))
	rows, err := p.db.QueryContext(ctx,
		`SELECT short, long FROM shorts WHERE user_id = $1`, uid)
	if err != nil {
		return nil, fmt.Errorf("error selecting rows: %w", err)
	}
	defer func() {
		if e := rows.Close(); e != nil {
			p.log.Warnf("error closing rows: %q", e.Error())
		}
	}()

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

// Delete deletes entries
func (p *Postgres) Delete(ctx context.Context, req entities.DeleteRequest) error {
	values := make([]string, len(req))
	for i, r := range req {
		values[i] = fmt.Sprintf("( short = '%s' AND user_id = '%s' )", r.URL, r.UID)
	}
	cmd :=
		"UPDATE shorts SET is_deleted = true WHERE " + strings.Join(values, " OR ") + ";"
	_, err := p.db.ExecContext(ctx, cmd)
	if err != nil {
		return fmt.Errorf("error marking for deletion: %w", err)
	}
	return nil
}

// CountUsersAndUrls - counts users and URLs
func (p *Postgres) CountUsersAndUrls(ctx context.Context) (entities.Stats, error) {
	var res entities.Stats

	err := p.db.QueryRowContext(
		ctx, `SELECT COUNT(short) FROM shorts WHERE is_deleted = false`).Scan(&res.URLs)
	if err != nil {
		return res, fmt.Errorf("error selecting URLs: %w", err)
	}

	err = p.db.QueryRowContext(
		ctx, `SELECT COUNT(DISTINCT user_id) FROM shorts WHERE is_deleted = false`).Scan(&res.Users)
	if err != nil {
		return res, fmt.Errorf("error selecting users: %w", err)
	}

	return res, nil
}

func migrate(ctx context.Context, log Logger, db *sql.DB) error {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	if err := goose.UpContext(ctx, db, "migrations"); err != nil {
		return fmt.Errorf("failed to migrate: %w", err)
	}

	return nil
}
