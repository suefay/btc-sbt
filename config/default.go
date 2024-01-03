package config

import (
	"fmt"
	"os"
	"path"
)

var (
	// Config file defaults to the current directory
	DefaultConfigFileName = path.Join(".", "config.yaml")

	// DB path defaults to the user directory
	DefaultDBPath = path.Join(".btc-sbt")

	// Key store defaults to the current directory
	DefaultKeyStorePath = path.Join(".", "keystore.txt")

	// Listener address default value
	DefaultListenerAddr = "0.0.0.0:80"
)

func init() {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("failed to get the user home directory: %v", err)
		os.Exit(1)
	}

	DefaultDBPath = path.Join(userHomeDir, DefaultDBPath)
}
