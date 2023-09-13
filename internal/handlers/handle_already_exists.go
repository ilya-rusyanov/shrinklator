package handlers

import (
	"errors"
	"net/http"

	"github.com/ilya-rusyanov/shrinklator/internal/storage"
)

func handleAlreadyExists(err error, statusCode *int) (value string, e error) {
	var exists storage.ErrAlreadyExists
	if errors.As(err, &exists) {
		*statusCode = http.StatusConflict
		return exists.StoredValue, nil
	}

	return "", err
}
