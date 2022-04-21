package getblock

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/egelis/GetBlock/internal/common"
)

const ErrGetBlockAPI = "getblock.io API error: %w"

type Client struct {
	APIKey string
}

func NewClient(key string) *Client {
	return &Client{APIKey: key}
}

func (c *Client) GetBlockByNumber(ctx context.Context, number string) (*common.Block, error) {
	body, err := c.postRequest(ctx, getBlockByNumberMethod, params{number, true})
	if err != nil {
		return nil, fmt.Errorf(ErrGetBlockAPI, err)
	}

	block := &common.Block{}
	if err = json.Unmarshal(body, &block); err != nil {
		return nil, err
	}

	return block, nil
}

func (c *Client) GetLastBlockNumber(ctx context.Context) (*common.LastBlockNum, error) {
	body, err := c.postRequest(ctx, blockNumberMethod, params{})
	if err != nil {
		return nil, fmt.Errorf(ErrGetBlockAPI, err)
	}

	lastBlockNum := &common.LastBlockNum{}
	if err = json.Unmarshal(body, &lastBlockNum); err != nil {
		return nil, fmt.Errorf(ErrGetBlockAPI, err)
	}

	return lastBlockNum, nil
}
