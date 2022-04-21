package getblock

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

const (
	getBlockByNumberMethod = "eth_getBlockByNumber"
	blockNumberMethod      = "eth_blockNumber"

	requestURL = "https://eth.getblock.io/mainnet/"
)

type params []interface{}

type RequestMessage struct {
	JsonRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	Id      string        `json:"id"`
}

func (c *Client) postRequest(ctx context.Context, method string, params params) ([]byte, error) {
	// Create request message
	bodyData := &RequestMessage{
		JsonRPC: "2.0",
		Method:  method,
		Params:  params,
		Id:      "getblock.io",
	}

	body, err := json.Marshal(bodyData)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, requestURL, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.APIKey)

	// Do request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	// Get body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}
