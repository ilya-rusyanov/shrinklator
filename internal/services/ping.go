package services

import (
	"context"
	"fmt"
)

// Database represents pingable DB object
type Database interface {
	Ping(context.Context) error
}

// Ping - is a usecase for pinging database
type Ping struct {
	db Database
}

// NewPing constructs Ping object
func NewPing(db Database) *Ping {
	return &Ping{
		db: db,
	}
}

// Ping pings database
func (p *Ping) Ping(ctx context.Context) error {
	err := p.db.Ping(ctx)

	if err != nil {
		return fmt.Errorf("DB ping failure: %w", err)
	}

	return nil
}
