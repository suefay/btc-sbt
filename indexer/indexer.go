package indexer

import (
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"

	"btc-sbt/config"
	"btc-sbt/logger"
	"btc-sbt/params"
	"btc-sbt/protocol"
	"btc-sbt/stacks/client/rpcclient"
	"btc-sbt/statemachine"
	"btc-sbt/store"
)

// Indexer defines the indexer struct
type Indexer struct {
	Client *rpcclient.Client // btc rpc client

	NetParams *chaincfg.Params // net params
	Params    *params.Params   // protocol params

	Parser       *protocol.Parser           // protocol parser
	StateMachine *statemachine.StateMachine // state machine

	Logger *logrus.Logger // logger

	interval time.Duration // block scanning interval
	done     bool          // indicates if the indexer is done
	stopped  chan struct{} // receive stop signal
	mu       sync.Mutex    // lock

	lastBlockHeight int64           // the height of the last indexed block
	lastBlockHash   *chainhash.Hash // the hash of the last indexed block
}

// NewIndexer creates a new Indexer instance
func NewIndexer(config *config.Config) (*Indexer, error) {
	client, err := rpcclient.NewClient(config.NodeRPCUrl, config.NodeRPCUser, config.NodeRPCPass)
	if err != nil {
		return nil, err
	}

	netParams := new(chaincfg.Params)
	protoParams := new(params.Params)

	switch config.NetVersion {
	case 1:
		netParams = &chaincfg.TestNet3Params
		protoParams = &params.TestNetParams
	case 2:
		netParams = &chaincfg.SigNetParams
		protoParams = &params.SigNetParams
	default:
		netParams = &chaincfg.MainNetParams
		protoParams = &params.MainNetParams
	}

	logger.Logger.SetLevel(logrus.Level(config.LogLevel))

	parser := protocol.NewParser(netParams)

	store, err := store.NewStore(config.DBPath)
	if err != nil {
		return nil, err
	}

	sm := statemachine.NewStateMachine(store, netParams, logger.Logger)

	indexer := &Indexer{
		Client:       client,
		NetParams:    netParams,
		Params:       protoParams,
		Parser:       parser,
		StateMachine: sm,
		Logger:       logger.Logger,
		interval:     config.IndexerInterval,
		stopped:      make(chan struct{}),
	}

	if err := indexer.loadStatus(); err != nil {
		return nil, err
	}

	return indexer, nil
}

// Start starts the indexer
func (i *Indexer) Start() error {
	i.onStart()

	go func() {
		i.startScanner()
		i.waitForStop()
	}()

	return nil
}

// Stop stops the running indexer
func (i *Indexer) Stop() error {
	i.onStop()

	return nil
}

// Stopped returns true if the indexer is stopped, false otherwise
func (i *Indexer) Stopped() bool {
	i.mu.Lock()
	defer i.mu.Unlock()

	return i.done
}

// loadStatus loads the indexing status
func (i *Indexer) loadStatus() error {
	lastBlockHeight, err := i.GetLastBlockHeight()
	if err != nil {
		return err
	}

	lastBlockHash, err := i.GetLastBlockHash()
	if err != nil {
		return err
	}

	i.lastBlockHeight = lastBlockHeight
	i.lastBlockHash = lastBlockHash

	return nil
}

// onStart is responsible for initial handling when started
func (i *Indexer) onStart() {
	i.Logger.Infof("indexer started")
	i.Logger.Infof("protocol params: %+v", *i.Params)

	if i.lastBlockHeight > 0 {
		i.Logger.Infof("last indexed block height: %d, hash: %s", i.lastBlockHeight, i.lastBlockHash)
	}

	latestBlockHeight, err := i.GetLatestBlockHeight()
	if err != nil {
		i.Logger.Errorf("failed to get the lastest block height: %v", err)
		return
	}

	i.Logger.Infof("latest block height: %d", latestBlockHeight)
}

// onStop is responsible to set the status that indicates the indexer is done
func (i *Indexer) onStop() {
	i.mu.Lock()
	i.done = true
	i.mu.Unlock()

	i.stopped <- struct{}{}
}

// waitForStop waits for the stop signal
func (i *Indexer) waitForStop() {
	<-i.stopped

	i.Logger.Info("indexer stopped")
}
