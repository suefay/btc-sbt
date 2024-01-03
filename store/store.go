package store

import (
	"bytes"
	"encoding/binary"

	"github.com/sirupsen/logrus"

	"github.com/cockroachdb/pebble"

	"btc-sbt/logger"
)

// Store defines a struct for data store
type Store struct {
	db *pebble.DB
}

// NewStore constructs a new Store instance
func NewStore(path string) (*Store, error) {
	db, err := pebble.Open(path, getDefaultOptions())
	if err != nil {
		return nil, err
	}

	return &Store{
		db: db,
	}, nil
}

// Set writes the given key-value into the store
func (s *Store) Set(key, value []byte) error {
	return s.db.Set(key, value, pebble.Sync)
}

// SetUint64 is a convenience to store the uint64 typed value
func (s *Store) SetUint64(key []byte, value uint64) error {
	bz := make([]byte, 8)
	binary.LittleEndian.PutUint64(bz, value)

	return s.Set(key, bz)
}

// SetInt64 is a convenience to store the int64 typed value
func (s *Store) SetInt64(key []byte, value int64) error {
	return s.SetUint64(key, uint64(value))
}

// Get retrieves the value of the given key
func (s *Store) Get(key []byte) ([]byte, error) {
	value, closer, err := s.db.Get(key)
	if err != nil {
		return nil, err
	}

	defer closer.Close()

	return value, nil
}

// GetUint64 is a convenience to get the uint64 typed value
func (s *Store) GetUint64(key []byte) (uint64, error) {
	value, err := s.Get(key)
	if err != nil {
		return 0, err
	}

	return binary.LittleEndian.Uint64(value), nil
}

// GetInt64 is a convenience to get the int64 typed value
func (s *Store) GetInt64(key []byte) (int64, error) {
	value, err := s.GetUint64(key)
	return int64(value), err
}

// Delete deletes the value by the given key
func (s *Store) Delete(key []byte) error {
	err := s.db.Delete(key, nil)
	if err != nil {
		return err
	}

	return nil
}

// Exist checks if the given key exists
func (s *Store) Exist(key []byte) (bool, error) {
	iter, err := s.db.NewIter(nil)
	if err != nil {
		return false, err
	}

	defer iter.Close()

	ok := iter.SeekGE(key)
	if !ok {
		return false, nil
	}

	return bytes.Equal(iter.Key(), key), nil
}

// getDefaultOptions gets the default options with the customized logger.
// The level is set to warn to silent info
func getDefaultOptions() *pebble.Options {
	logger.Logger.SetLevel(logrus.WarnLevel)

	return &pebble.Options{
		Logger: logger.Logger,
	}
}
