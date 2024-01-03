package initiator

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"

	"btc-sbt/config"
	"btc-sbt/logger"
	"btc-sbt/stacks/client/base"
	"btc-sbt/stacks/client/unisat"
)

// Initiator defines the initiator struct which is intended to initiate the protocol operations
type Initiator struct {
	RPCClient *rpcclient.Client // btc rpc client

	UnisatClient *unisat.Client // unisat client

	NetParams *chaincfg.Params // net params
	Config    *config.Config   // config

	Logger *logrus.Logger // logger
}

// NewInitiator creates a new Initiator instance
func NewInitiator(config *config.Config) (*Initiator, error) {
	var netParams *chaincfg.Params

	switch config.NetVersion {
	case 0:
		netParams = &chaincfg.MainNetParams
	case 1:
		netParams = &chaincfg.TestNet3Params
	default:
		netParams = &chaincfg.SigNetParams
	}

	rpcClient, err := createRPCClient(config, netParams)
	if err != nil {
		return nil, err
	}

	unisatClient := createUnisatClient(config)
	if unisatClient == nil {
		return nil, fmt.Errorf("unisat client not created")
	}

	logger.Logger.SetLevel(logrus.Level(config.LogLevel))

	return &Initiator{
		RPCClient:    rpcClient,
		UnisatClient: unisatClient,
		NetParams:    netParams,
		Config:       config,
		Logger:       logger.Logger,
	}, nil
}

// createRPCClient creates the rpc client
func createRPCClient(config *config.Config, netParams *chaincfg.Params) (*rpcclient.Client, error) {
	connCfg := &rpcclient.ConnConfig{
		Host:         config.NodeRPCUrl,
		User:         config.NodeRPCUser,
		Pass:         config.NodeRPCPass,
		HTTPPostMode: true,
		DisableTLS:   true,
	}

	rpcClient, err := rpcclient.New(connCfg, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create rpc client, err: %v", err)
	}

	_, err = rpcClient.GetBlockChainInfo()
	if err != nil {
		return nil, fmt.Errorf("failed to get blockchain info, err: %v", err)
	}

	return rpcClient, nil
}

// createUnisatClient creates a unisat client
func createUnisatClient(config *config.Config) *unisat.Client {
	if len(config.UnisatAPI) > 0 {
		return unisat.NewClient(
			config.UnisatAPI,
			base.NewClient(config.Retries+1, config.Interval),
		)
	}

	return nil
}
