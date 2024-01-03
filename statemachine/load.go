package statemachine

import (
	"github.com/btcsuite/btcd/chaincfg/chainhash"

	"btc-sbt/store"
	"btc-sbt/types"
)

// GetAllSBTs queries all the SBTs from the store
func (sm *StateMachine) GetAllSBTs() ([]*types.SBTs, error) {
	iter, err := sm.Store.Iterator(GetSBTsKeyPrefix())
	if err != nil {
		return nil, err
	}

	defer iter.Close()

	collections := make([]*types.SBTs, 0)

	for iter.First(); iter.Valid(); iter.Next() {
		var sbts types.SBTs
		if err := sbts.Unmarshal(iter.Value()); err != nil {
			return nil, err
		}

		supply, err := sm.GetSBTsSupply(sbts.Symbol)
		if err != nil {
			return nil, err
		}

		sbts.TotalSupply = supply

		collections = append(collections, &sbts)
	}

	return collections, nil
}

// GetSBTs queries the SBTs by the given symbol from the store
func (sm *StateMachine) GetSBTs(symbol string) (*types.SBTs, error) {
	key := GetSBTsKey(symbol)

	bz, err := sm.Store.Get(key)
	if err != nil && !store.IsNotFoundErr(err) {
		return nil, err
	}

	if len(bz) == 0 {
		return nil, nil
	}

	var sbts types.SBTs
	if err := sbts.Unmarshal(bz); err != nil {
		return nil, err
	}

	supply, err := sm.GetSBTsSupply(symbol)
	if err != nil {
		return nil, err
	}

	sbts.TotalSupply = supply

	return &sbts, nil
}

// GetSBT queries the SBT token by the given symbol and token id from the store
func (sm *StateMachine) GetSBT(symbol string, id uint64) (*types.SBT, error) {
	key := GetSBTKey(symbol, id)

	bz, err := sm.Store.Get(key)
	if err != nil && !store.IsNotFoundErr(err) {
		return nil, err
	}

	if len(bz) == 0 {
		return nil, nil
	}

	var sbt types.SBT
	if err := sbt.Unmarshal(bz); err != nil {
		return nil, err
	}

	return &sbt, nil
}

// GetSBTsSequence queries the current SBTs sequence from the store
func (sm *StateMachine) GetSBTsSequence() (uint64, error) {
	key := GetSBTsSequenceKey()

	seq, err := sm.Store.GetUint64(key)
	if err != nil && !store.IsNotFoundErr(err) {
		return 0, err
	}

	return seq, nil
}

// GetSBTsSupply queries the current supply of the given SBTs from the store
func (sm *StateMachine) GetSBTsSupply(symbol string) (uint64, error) {
	key := GetSBTsSupplyKey(symbol)

	supply, err := sm.Store.GetUint64(key)
	if err != nil && !store.IsNotFoundErr(err) {
		return 0, err
	}

	return supply, nil
}

// GetOwnedSBTs queries the SBT tokens owned by the given owner from the store
func (sm *StateMachine) GetOwnedSBTs(owner string) ([]*types.CompactSBT, error) {
	iter, err := sm.Store.Iterator(GetOwnerSBTKeyPrefix(owner))
	if err != nil {
		return nil, err
	}

	defer iter.Close()

	sbts := make([]*types.CompactSBT, 0)

	for iter.First(); iter.Valid(); iter.Next() {
		var sbt types.CompactSBT
		if err := sbt.Unmarshal(iter.Value()); err != nil {
			return nil, err
		}

		sbts = append(sbts, &sbt)
	}

	return sbts, nil
}

// GetOwnedSBT queries the specified SBT token owned by the given owner from the store
func (sm *StateMachine) GetOwnedSBT(owner string, symbol string) (*types.CompactSBT, error) {
	key := GetOwnerSBTKey(owner, symbol)

	bz, err := sm.Store.Get(key)
	if err != nil && !store.IsNotFoundErr(err) {
		return nil, err
	}

	if len(bz) == 0 {
		return nil, nil
	}

	var sbt types.CompactSBT
	if err := sbt.Unmarshal(bz); err != nil {
		return nil, err
	}

	return &sbt, nil
}

// GetLastBlockHeight queries the last block height of the indexer from the store
func (sm *StateMachine) GetLastBlockHeight() (int64, error) {
	key := GetIndexerLastBlockHeightKey()

	blockHeight, err := sm.Store.GetInt64(key)
	if err != nil && !store.IsNotFoundErr(err) {
		return 0, err
	}

	return blockHeight, nil
}

// GetLastBlockHash queries the last block hash of the indexer from the store
func (sm *StateMachine) GetLastBlockHash() (*chainhash.Hash, error) {
	key := GetIndexerLastBlockHashKey()

	blockHash, err := sm.Store.Get(key)
	if err != nil && !store.IsNotFoundErr(err) {
		return nil, err
	}

	if blockHash == nil {
		return nil, nil
	}

	return chainhash.NewHash(blockHash)
}

// SBTsExists returns true if the given SBTs exists, false otherwise
func (sm *StateMachine) SBTsExists(symbol string) (bool, error) {
	return sm.Store.Exist(GetSBTsKey(symbol))
}

// SBTExists returns true if the given SBT token exist, false otherwise
func (sm *StateMachine) SBTExists(symbol string, id uint64) (bool, error) {
	return sm.Store.Exist(GetSBTKey(symbol, id))
}

// HasOwnedSBT returns true if the given address has owned the SBT, false otherwise
func (sm *StateMachine) HasOwnedSBT(address string, symbol string) (bool, error) {
	return sm.Store.Exist(GetOwnerSBTKey(address, symbol))
}
