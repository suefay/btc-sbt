package statemachine

import (
	"github.com/sirupsen/logrus"

	"github.com/btcsuite/btcd/chaincfg"

	"btc-sbt/store"
)

// StateMachine handles the state transition of the BTC-SBT protocol
type StateMachine struct {
	Store *store.Store // store

	NetParams *chaincfg.Params // net params

	Logger *logrus.Logger // logger
}

// NewStateMachine creates a new StateMachine instance
func NewStateMachine(store *store.Store, netParams *chaincfg.Params, logger *logrus.Logger) *StateMachine {
	return &StateMachine{
		Store:     store,
		NetParams: netParams,
		Logger:    logger,
	}
}
