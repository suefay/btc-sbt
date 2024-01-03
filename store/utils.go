package store

import (
	"github.com/cockroachdb/pebble"
)

// IsNotFoundErr returns true if the given error is ErrNotFound, false otherwise
func IsNotFoundErr(err error) bool {
	return err == pebble.ErrNotFound
}
