package storage

import (
	"context"
	"fmt"
)

type NoDB struct {
}

func NewNoDB() *NoDB {
	return &NoDB{}
}

func (d *NoDB) Ping(context.Context) error {
	return fmt.Errorf("DB is not available")
}
