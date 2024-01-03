package rpcclient

import (
	"fmt"

	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/wire"
)

// Client wraps the rpcclient.Client
type Client struct {
	inner *rpcclient.Client
}

// NewClient creates a RPC client instance
func NewClient(rpcUrl, rpcUser, rpcPass string) (*Client, error) {
	connCfg := &rpcclient.ConnConfig{
		Host:         rpcUrl,
		User:         rpcUser,
		Pass:         rpcPass,
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

	return &Client{inner: rpcClient}, nil
}

// GetLatestBlockHeight gets the height of the lastest block
func (c *Client) GetLatestBlockHeight() (int64, error) {
	height, err := c.inner.GetBlockCount()

	return height, err
}

// GetLatestBlock gets the lastest block
func (c *Client) GetLatestBlock() (*wire.MsgBlock, error) {
	hash, err := c.inner.GetBestBlockHash()
	if err != nil {
		return nil, err
	}

	return c.inner.GetBlock(hash)
}

// GetBlock gets the block by the given height
func (c *Client) GetBlock(height int64) (*wire.MsgBlock, error) {
	hash, err := c.inner.GetBlockHash(height)
	if err != nil {
		return nil, err
	}

	return c.inner.GetBlock(hash)
}
