package services

import (
	"context"
	"fmt"
)

type Database interface {
	Ping(context.Context) error
}

type Ping struct {
	db Database
}

func NewPing(db Database) *Ping {
	return &Ping{
		db: db,
	}
}

func (p *Ping) Ping(ctx context.Context) error {
	err := p.db.Ping(ctx)

	if err != nil {
		return fmt.Errorf("DB ping failure: %w", err)
	}

	return nil
}
