package statemachine

import (
	"encoding/binary"
	"strings"
)

var (
	KEY_SEPARATOR = byte(0x00)

	SBTS_KEY_PREFIX = []byte{0x01}
	SBT_KEY_PREFIX  = []byte{0x02}

	SBTS_SEQUENCE_KEY_PREFIX = []byte{0x03}
	SBTS_SUPPLY_KEY_PREFIX   = []byte{0x04}

	OWNER_SBT_KEY_PREFIX = []byte{0x05}

	INDEXER_STATUS_LAST_BLOCK_HEIGHT_KEY = []byte{0x06}
	INDEXER_STATUS_LAST_BLOCK_HASH_KEY   = []byte{0x07}
)

// GetSBTsKeyPrefix gets the key prefix for iteration over all the SBTs
func GetSBTsKeyPrefix() []byte {
	return SBTS_KEY_PREFIX
}

// GetSBTsKey gets the store key for the SBTs by the given symbol
func GetSBTsKey(symbol string) []byte {
	return append(SBTS_KEY_PREFIX, []byte(strings.ToLower(symbol))...)
}

// GetSBTKey gets the store key for the SBT token by the given symbol and id
func GetSBTKey(symbol string, id uint64) []byte {
	idBz := make([]byte, 8)
	binary.BigEndian.PutUint64(idBz, id)

	return append(append(SBT_KEY_PREFIX, []byte(strings.ToLower(symbol))...), idBz...)
}

// GetSBTsSequenceKey gets the store key for the current SBTs sequence
func GetSBTsSequenceKey() []byte {
	return SBTS_SEQUENCE_KEY_PREFIX
}

// GetSBTsSupplyKey gets the store key for the current supply of the given SBTs
func GetSBTsSupplyKey(symbol string) []byte {
	return append(SBTS_SUPPLY_KEY_PREFIX, []byte(strings.ToLower(symbol))...)
}

// GetOwnerSBTKey gets the store key for the SBT token owned by the given owner
func GetOwnerSBTKey(owner string, symbol string) []byte {
	prefix := GetOwnerSBTKeyPrefix(strings.ToLower(owner))

	return append(prefix, []byte(strings.ToLower(symbol))...)
}

// GetOwnerSBTKeyPrefix gets the key prefix for iteration over the SBT tokens owned by the given owner
func GetOwnerSBTKeyPrefix(owner string) []byte {
	prefix := append(OWNER_SBT_KEY_PREFIX, []byte(strings.ToLower(owner))...)
	prefix = append(prefix, KEY_SEPARATOR)

	return prefix
}

// GetIndexerLastBlockHeightKey gets the store key for the last block height of the indexer
func GetIndexerLastBlockHeightKey() []byte {
	return INDEXER_STATUS_LAST_BLOCK_HEIGHT_KEY
}

// GetIndexerLastBlockHashKey gets the store key for the last block hash of the indexer
func GetIndexerLastBlockHashKey() []byte {
	return INDEXER_STATUS_LAST_BLOCK_HASH_KEY
}
