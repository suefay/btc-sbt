package node

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/sirupsen/logrus"

	"btc-sbt/config"
	"btc-sbt/indexer"
	"btc-sbt/logger"
	"btc-sbt/server"
)

// Node represents the node instance for the BTC-SBT protocol
type Node struct {
	Indexer *indexer.Indexer // indexer

	APIService *server.APIService // api service

	Config *config.Config // config

	Logger *logrus.Logger // logger
}

// CreateNode creates a new Node instance
func CreateNode(config *config.Config) (*Node, error) {
	indexer, err := indexer.NewIndexer(config)
	if err != nil {
		return nil, err
	}

	logger.Logger.SetLevel(logrus.Level(config.LogLevel))

	apiService := server.NewAPIService(indexer, logger.Logger)

	return &Node{
		Indexer:    indexer,
		APIService: apiService,
		Config:     config,
		Logger:     logger.Logger,
	}, nil
}

// Start starts the node
func (n *Node) Start() {
	n.Logger.Infof("starting node")

	n.startHandler()

	n.waitForStop()
}

// waitForStop waits for the stop signal
func (n *Node) waitForStop() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig

	n.Logger.Infof("Stopping the node")

	if err := n.stopHandler(); err != nil {
		n.Logger.Errorf("%v", err)
	}
}

// startHandler handles the start logic
func (n *Node) startHandler() {
	if err := n.Indexer.Start(); err != nil {
		n.Logger.Fatalf("failed to start the indexer: %v", err)
	}

	go func() {
		if err := n.APIService.Start(n.Config.ListenerAddr); err != nil {
			n.Logger.Fatalf("failed to start the API service: %v", err)
		}
	}()
}

// stopHandler handles the stop logic
func (n *Node) stopHandler() error {
	if err := n.Indexer.Stop(); err != nil {
		return fmt.Errorf("failed to stop the indexer: %v", err)
	}

	if err := n.APIService.Stop(); err != nil {
		return fmt.Errorf("failed to stop the API service: %v", err)
	}

	return nil
}
