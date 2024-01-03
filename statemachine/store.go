package statemachine

import (
	"github.com/btcsuite/btcd/chaincfg/chainhash"

	"btc-sbt/types"
)

// SetSBTs sets the given SBTs in the store
func (sm *StateMachine) SetSBTs(sbts *types.SBTs) error {
	bz, err := sbts.Marshal()
	if err != nil {
		return err
	}

	key := GetSBTsKey(sbts.Symbol)

	return sm.Store.Set(key, bz)
}

// SetSBT sets the given SBT token in the store
func (sm *StateMachine) SetSBT(sbt *types.SBT) error {
	bz, err := sbt.Marshal()
	if err != nil {
		return err
	}

	key := GetSBTKey(sbt.Symbol, sbt.Id)

	return sm.Store.Set(key, bz)
}

// SetSBTsSequence sets the current SBTs sequence in the store
func (sm *StateMachine) SetSBTsSequence(sequence uint64) error {
	key := GetSBTsSequenceKey()

	return sm.Store.SetUint64(key, sequence)
}

// IncreaseSBTsSequence increases the current SBTs sequence by 1 and returns the new sequence
func (sm *StateMachine) IncreaseSBTsSequence() (uint64, error) {
	seq, err := sm.GetSBTsSequence()
	if err != nil {
		return 0, err
	}

	if err := sm.SetSBTsSequence(seq + 1); err != nil {
		return 0, err
	}

	return seq + 1, nil
}

// SetSBTsSupply sets the current supply of the given SBTs in the store
func (sm *StateMachine) SetSBTsSupply(symbol string, supply uint64) error {
	key := GetSBTsSupplyKey(symbol)

	return sm.Store.SetUint64(key, supply)
}

// SetOwnerSBT sets the SBT token by the given owner in the store
func (sm *StateMachine) SetOwnerSBT(owner string, sbt *types.SBT) error {
	bz, err := sbt.Compact().Marshal()
	if err != nil {
		return err
	}

	key := GetOwnerSBTKey(owner, sbt.Symbol)

	return sm.Store.Set(key, bz)
}

// SetLastBlockHeight sets the last block height of the indexer in the store
func (sm *StateMachine) SetLastBlockHeight(height int64) error {
	key := GetIndexerLastBlockHeightKey()

	return sm.Store.SetInt64(key, height)
}

// SetLastBlockHash sets the last block hash of the indexer in the store
func (sm *StateMachine) SetLastBlockHash(hash chainhash.Hash) error {
	key := GetIndexerLastBlockHashKey()

	return sm.Store.Set(key, hash[:])
}
