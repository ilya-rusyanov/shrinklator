package grpcsrv

import (
	"errors"

	"github.com/ilya-rusyanov/shrinklator/internal/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func handleAlreadyExists(err error, statusCode *error) (value string, e error) {
	var exists storage.ErrAlreadyExists

	if errors.As(err, &exists) {
		*statusCode = status.Errorf(codes.AlreadyExists, "given URL already exists")
		return exists.StoredValue, nil
	}

	return "", err
}
