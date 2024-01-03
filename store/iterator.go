package store

import (
	"github.com/cockroachdb/pebble"
)

// Iterator gets the iterator by the given key prefix
func (s *Store) Iterator(prefix []byte) (*pebble.Iterator, error) {
	return s.db.NewIter(getPrefixIterOptions(prefix))
}

// getPrefixIterOptions gets the iterator options by the given key prefix
func getPrefixIterOptions(prefix []byte) *pebble.IterOptions {
	return &pebble.IterOptions{
		LowerBound: prefix,
		UpperBound: getKeyUpperBound(prefix),
	}
}

// getKeyUpperBound returns the upper bound for iteration over the given prefix
func getKeyUpperBound(prefix []byte) []byte {
	end := make([]byte, len(prefix))
	copy(end, prefix)

	for i := len(end) - 1; i >= 0; i-- {
		end[i] = end[i] + 1
		if end[i] != 0 {
			return end[:i+1]
		}
	}

	return nil
}
