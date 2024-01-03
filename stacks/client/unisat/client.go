package unisat

import (
	"fmt"
	"net/http"

	"btc-sbt/stacks/basics"
	"btc-sbt/stacks/client/base"
)

// Client defines the unisat client
type Client struct {
	BaseClient *base.Client

	UnisatAPI string
}

// NewClient creates a unisat client instance
func NewClient(unisatAPI string, baseClient *base.Client) *Client {
	return &Client{
		BaseClient: baseClient,
		UnisatAPI:  unisatAPI,
	}
}

// GetAddressUtxos queries the non-ordinals utxos of the given address
func (c *Client) GetAddressUtxos(address string) ([]*basics.UTXO, error) {
	url := fmt.Sprintf("%s/address/btc-utxo?address=%s", c.UnisatAPI, address)

	opts := c.BaseClient.GetBaseOptions()
	opts.Headers["X-Client"] = UNISAT_WALLET_CLIENT
	opts.Headers["X-Version"] = UNISAT_WALLET_VERSION

	statusCode, resp, err := c.BaseClient.Request(http.MethodGet, url, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to query utxos, err: %v", err)
	}

	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to query utxos, status code: %d, response: %s", statusCode, string(resp))
	}

	var r GetUtxosResponse
	if err := r.UnmarshalJSON(resp); err != nil {
		return nil, fmt.Errorf("failed to query utxos: invalid response, err: %v", err)
	}

	return r.GetUtxos()
}
