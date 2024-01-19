package ordinals

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/btcsuite/btcd/wire"

	"btc-sbt/stacks/client/base"
)

const HTTP_HEADER_ACCEPT_JSON = "application/json"

// Client defines the ordinals client
type Client struct {
	BaseClient *base.Client

	IndexerAPI string
}

// NewClient creates an ordinals client instance
func NewClient(indexerAPI string, baseClient *base.Client) *Client {
	return &Client{
		BaseClient: baseClient,
		IndexerAPI: indexerAPI,
	}
}

// GetInscription gets the given inscription
func (c *Client) GetInscription(id string) (*Inscription, error) {
	url := fmt.Sprintf("%s/inscription/%s", c.IndexerAPI, id)

	opts := c.BaseClient.GetBaseOptions()
	opts.Headers["Accept"] = HTTP_HEADER_ACCEPT_JSON

	statusCode, resp, err := c.BaseClient.Request(http.MethodGet, url, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to query the inscription, err: %v", err)
	}

	switch statusCode {
	case http.StatusOK:
		var i Inscription
		if err := json.Unmarshal(resp, &i); err != nil {
			return nil, fmt.Errorf("failed to query the inscription: invalid response, err: %v", err)
		}

		return &i, nil

	case http.StatusNotFound:
		return nil, nil

	default:
		return nil, fmt.Errorf("failed to query the inscription, status code: %d, response: %s", statusCode, string(resp))
	}
}

// GetInscriptions retrieves all inscriptions of the given address
func (c *Client) GetInscriptions(address string) ([]string, error) {
	url := fmt.Sprintf("%s/address/%s", c.IndexerAPI, address)

	opts := c.BaseClient.GetBaseOptions()
	opts.Headers["Accept"] = HTTP_HEADER_ACCEPT_JSON

	statusCode, resp, err := c.BaseClient.Request(http.MethodGet, url, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to query inscriptions, err: %v", err)
	}

	switch statusCode {
	case http.StatusOK:
		var r GetInscriptionsResponse
		if err := r.UnmarshalJSON(resp); err != nil {
			return nil, fmt.Errorf("failed to query inscriptions: invalid response, err: %v", err)
		}

		return r.GetInscriptions(), nil

	case http.StatusNotFound:
		return nil, nil

	default:
		return nil, fmt.Errorf("failed to query inscriptions, status code: %d, response: %s", statusCode, string(resp))
	}
}

// GetInscriptionsByOutput gets inscriptions by the given output point
func (c *Client) GetInscriptionsByOutput(output *wire.OutPoint) ([]string, error) {
	url := fmt.Sprintf("%s/output/%s", c.IndexerAPI, output.String())

	opts := c.BaseClient.GetBaseOptions()
	opts.Headers["Accept"] = HTTP_HEADER_ACCEPT_JSON

	statusCode, resp, err := c.BaseClient.Request(http.MethodGet, url, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to query inscriptions by output, err: %v", err)
	}

	switch statusCode {
	case http.StatusOK:
		var r GetInscriptionsByOutputResponse
		if err := r.Unmarshal(resp); err != nil {
			return nil, fmt.Errorf("failed to query inscriptions by output: invalid response, err: %v", err)
		}

		return r.GetInscriptions(), nil

	case http.StatusNotFound:
		return nil, nil

	default:
		return nil, fmt.Errorf("failed to query inscriptions by output, status code: %d, response: %s", statusCode, string(resp))
	}
}
