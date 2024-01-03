package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Config defines the config
type Config struct {
	NodeRPCUrl  string // node rpc url
	NodeRPCUser string // node rpc username
	NodeRPCPass string // node rpc passphrase

	NetVersion uint8 // net version: 0(mainnet), 1(testnet), 2(signet)

	UnisatAPI string // unisat api

	IndexerInterval time.Duration // indexer interval

	DBPath string // db path

	KeyStorePath string // key store path

	FeeRate int64 // fee rate

	Retries  int           // retry count
	Interval time.Duration // retry interval

	ListenerAddr string // listener address for web server

	LogLevel uint32 // logging level
}

// NewConfig creates a new Config instance
func NewConfig(
	nodeRPCUrl,
	nodeRPCUser,
	nodeRPCPass string,
	netVersion uint8,
	unisatAPI string,
	indexerInterval time.Duration,
	dbPath,
	keyStorePath string,
	feeRate int64,
	retries int,
	interval time.Duration,
	listenerAddr string,
	logLevel uint32,
) *Config {
	return &Config{
		NodeRPCUrl:      nodeRPCUrl,
		NodeRPCUser:     nodeRPCUser,
		NodeRPCPass:     nodeRPCPass,
		NetVersion:      netVersion,
		UnisatAPI:       unisatAPI,
		IndexerInterval: indexerInterval,
		DBPath:          dbPath,
		KeyStorePath:    keyStorePath,
		FeeRate:         feeRate,
		Retries:         retries,
		Interval:        interval,
		ListenerAddr:    listenerAddr,
		LogLevel:        logLevel,
	}
}

// NewConfigFromViper creates a new Config instance from viper
func NewConfigFromViper(v *viper.Viper) (*Config, error) {
	nodeRPCUrl := v.GetString("node.rpc_url")
	nodeRPCUser := v.GetString("node.rpc_user")
	nodeRPCPass := v.GetString("node.rpc_pass")

	netVersion := uint8(v.GetUint("node.net_version"))
	if netVersion > 2 {
		return nil, fmt.Errorf("invalid net version: only 0, 1 or 2 allowed, %d given", netVersion)
	}

	unisatAPI := v.GetString("unisat.api")

	indexerInterval := v.GetDuration("indexer.interval")

	dbPath := v.GetString("db.path")
	if len(dbPath) == 0 {
		dbPath = DefaultDBPath
	}

	keyStorePath := v.GetString("key_store.path")
	if len(keyStorePath) == 0 {
		keyStorePath = DefaultKeyStorePath
	}

	feeRate := v.GetInt64("fee_rate")

	retries := v.GetInt("general.retries")
	interval := v.GetDuration("general.interval")

	listenerAddr := v.GetString("server.listener_address")
	if len(listenerAddr) == 0 {
		listenerAddr = DefaultListenerAddr
	}

	logLevel := v.GetUint32("log.level")

	return NewConfig(
		nodeRPCUrl,
		nodeRPCUser,
		nodeRPCPass,
		netVersion,
		unisatAPI,
		indexerInterval,
		dbPath,
		keyStorePath,
		feeRate,
		retries,
		interval,
		listenerAddr,
		logLevel,
	), nil
}
