package service

import (
	"context"
	"fmt"
	"github.com/egelis/GetBlock/internal/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

type MockAPI struct {
	TestBlocks map[string]*common.Block
}

func (m *MockAPI) GetLastBlockNumber(_ context.Context) (*common.LastBlockNum, error) {
	lastBlock := len(m.TestBlocks) - 1
	return &common.LastBlockNum{Number: fmt.Sprintf("0x%d", lastBlock)}, nil
}

func (m *MockAPI) GetBlockByNumber(_ context.Context, number string) (*common.Block, error) {
	if block, ok := m.TestBlocks[number]; ok {
		return block, nil
	}

	return &common.Block{
		BlockContent: &common.BlockContent{
			Transactions: []*common.Transaction{},
		},
	}, nil
}

func TestService_GetAddrMostChangedBalance(t *testing.T) {
	testTable := []struct {
		name       string
		testBlocks map[string]*common.Block
		expected   common.MostChanged
	}{
		{
			name:       "The target balance is in one block",
			testBlocks: TargetAddrInOneBlock,
			expected: common.MostChanged{
				Address: "0x111",
				Value:   "0x4",
			},
		},
		{
			name:       "The target balance is in several blocks",
			testBlocks: TargetAddrInSeveralBlocks,
			expected: common.MostChanged{
				Address: "0x444",
				Value:   "0x3",
			},
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			testAPI := &MockAPI{TestBlocks: test.testBlocks}
			srv := NewService(testAPI)

			res, err := srv.GetAddrMostChangedBalance(context.Background())

			assert.Equal(t, nil, err)
			assert.Equal(t, test.expected, *res)
		})
	}
}

var (
	TargetAddrInOneBlock = map[string]*common.Block{
		"0x1": {
			BlockContent: &common.BlockContent{
				Transactions: []*common.Transaction{
					{
						From:  "0x111",
						To:    "0x2",
						Value: "0x1",
					},
					{
						From:  "0x111",
						To:    "0x3",
						Value: "0x1",
					},
					{
						From:  "0x111",
						To:    "0x3",
						Value: "0x1",
					},
					{
						From:  "0x111",
						To:    "0x4",
						Value: "0x1",
					},
				},
			},
		},
		"0x0": {
			BlockContent: &common.BlockContent{
				Transactions: []*common.Transaction{
					{
						From:  "0x4",
						To:    "0x5",
						Value: "0x1",
					},
					{
						From:  "0x5",
						To:    "0x3",
						Value: "0x1",
					},
				},
			},
		},
	}

	//

	TargetAddrInSeveralBlocks = map[string]*common.Block{
		"0x2": {
			BlockContent: &common.BlockContent{
				Transactions: []*common.Transaction{
					{
						From:  "0x1",
						To:    "0x2",
						Value: "0x1",
					},
					{
						From:  "0x1",
						To:    "0x3",
						Value: "0x1",
					},
					{
						From:  "0x3",
						To:    "0x1",
						Value: "0x1",
					},
					{
						From:  "0x3",
						To:    "0x444",
						Value: "0x1",
					},
				},
			},
		},
		"0x1": {
			BlockContent: &common.BlockContent{
				Transactions: []*common.Transaction{
					{
						From:  "0x444",
						To:    "0x3",
						Value: "0x1",
					},
					{
						From:  "0x5",
						To:    "0x444",
						Value: "0x1",
					},
					{
						From:  "0x6",
						To:    "0x444",
						Value: "0x1",
					},
				},
			},
		},
		"0x0": {
			BlockContent: &common.BlockContent{
				Transactions: []*common.Transaction{
					{
						From:  "0x7",
						To:    "0x8",
						Value: "0x1",
					},
					{
						From:  "0x8",
						To:    "0x444",
						Value: "0x1",
					},
				},
			},
		},
	}
)
