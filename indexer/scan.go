package indexer

import (
	"time"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
)

// startScanner starts to scan blocks
func (i *Indexer) startScanner() {
	for {
		i.scanBlocks()

		if i.Stopped() {
			return
		}

		time.Sleep(i.interval)
	}
}

// scanBlocks scans blocks
func (i *Indexer) scanBlocks() {
	currentHeight, err := i.GetLatestBlockHeight()
	if err != nil {
		i.Logger.Errorf("failed to retrieve the latest block height: %v", err)
		return
	}

	if i.lastBlockHeight < i.Params.ActivationBlockHeight {
		i.lastBlockHeight = i.Params.ActivationBlockHeight - 1
	}

	if currentHeight <= i.lastBlockHeight {
		return
	}

	i.scanBlocksByRange(i.lastBlockHeight+1, currentHeight)
}

// scanBlocksByRange scans blocks by the given range
func (i *Indexer) scanBlocksByRange(startHeight int64, endHeight int64) {
	for h := startHeight; h <= endHeight && !i.Stopped(); {
		block, err := i.Client.GetBlock(h)
		if err != nil {
			i.Logger.Errorf("failed to retrieve the block, height: %d, err: %v", h, err)
			continue
		}

		i.onBlock(h, block)

		h++
	}
}

// onBlock handles the given block
func (i *Indexer) onBlock(height int64, block *wire.MsgBlock) {
	if i.isReorged(block) {
		i.Logger.Fatalf("chain reorg detected, please reindex. The previous block hash of the indexing block: %s, last indexed block hash: %s", block.Header.PrevBlock, i.lastBlockHash)
	}

	if err := i.parseBTCSBTProtocol(height, block); err != nil {
		i.Logger.Fatalf("parsing failed, block height: %d, block hash: %s, err: %v", height, block.BlockHash(), err)
	}

	if err := i.saveStatus(height, block.BlockHash()); err != nil {
		i.Logger.Fatalf("failed to save the indexing status: %v", err)
	}

	i.Logger.Infof("block indexed: %d", height)
}

// saveStatus saves the current indexing status
func (i *Indexer) saveStatus(blockHeight int64, blockHash chainhash.Hash) error {
	// store status

	if err := i.StateMachine.SetLastBlockHeight(blockHeight); err != nil {
		return err
	}

	if err := i.StateMachine.SetLastBlockHash(blockHash); err != nil {
		return err
	}

	// in-memory status

	i.lastBlockHeight = blockHeight
	i.lastBlockHash = &blockHash

	return nil
}

// isReorged returns true if the given block does not point to the last indexed block, false otherwise
func (i *Indexer) isReorged(block *wire.MsgBlock) bool {
	return i.lastBlockHash != nil && !block.Header.PrevBlock.IsEqual(i.lastBlockHash)
}
