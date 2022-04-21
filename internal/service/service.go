package service

import (
	"context"
	"fmt"
	"github.com/egelis/GetBlock/internal/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"golang.org/x/sync/errgroup"
	"math/big"
	"sync"
)

const numOfCalcBlocks = 100

const (
	ErrBadValueTxs      = "got bad value in transactions: %w"
	ErrBadValueBlockNum = "got bad value in block number: %w"
)

type ClientAPI interface {
	GetBlockByNumber(ctx context.Context, number string) (*common.Block, error)
	GetLastBlockNumber(ctx context.Context) (*common.LastBlockNum, error)
}

type Service struct {
	clientAPI ClientAPI
}

func NewService(client ClientAPI) *Service {
	return &Service{
		clientAPI: client,
	}
}

func (s *Service) GetAddrMostChangedBalance(ctx context.Context) (*common.MostChanged, error) {
	eg, ctx := errgroup.WithContext(ctx)

	addresses := map[string]*big.Int{}
	blockNums := make(chan string, numOfCalcBlocks)

	// Calculating block numbers
	eg.Go(func() error {
		defer close(blockNums)

		blockNumStr, err := s.clientAPI.GetLastBlockNumber(ctx)
		if err != nil {
			return err
		}

		blockNum, err := hexutil.DecodeBig(blockNumStr.Number)
		if err != nil {
			return fmt.Errorf(ErrBadValueBlockNum, err)
		}

		decrement := big.NewInt(1)
		for i := 0; i < numOfCalcBlocks; i++ {
			blockNums <- hexutil.EncodeBig(blockNum)
			blockNum.Sub(blockNum, decrement)
		}

		return nil
	})

	m := &sync.Mutex{}

	// Get addresses from transactions of all blocks
	for blockNumber := range blockNums {
		blockNumber := blockNumber

		eg.Go(func() error {
			block, err := s.clientAPI.GetBlockByNumber(ctx, blockNumber)
			if err != nil {
				return err
			}

			if block.Transactions != nil {
				if err := s.getAllAddrFromBlockTxs(m, addresses, block.Transactions); err != nil {
					return err
				}
			}

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	return s.findMaxAddr(addresses), nil
}

func (s *Service) getAllAddrFromBlockTxs(m *sync.Mutex, addresses map[string]*big.Int,
	transactions []*common.Transaction) error {
	for _, tr := range transactions {
		value, err := hexutil.DecodeBig(tr.Value)
		if err != nil {
			return fmt.Errorf(ErrBadValueTxs, err)
		}

		m.Lock()
		if _, ok := addresses[tr.From]; !ok {
			newValue := new(big.Int).Set(value)
			addresses[tr.From] = newValue.Neg(newValue)
		} else {
			addresses[tr.From].Sub(addresses[tr.From], value)
		}

		if _, ok := addresses[tr.To]; !ok {
			newValue := new(big.Int).Set(value)
			addresses[tr.To] = newValue
		} else {
			addresses[tr.To].Add(addresses[tr.To], value)
		}
		m.Unlock()
	}

	return nil
}

func (s *Service) findMaxAddr(addresses map[string]*big.Int) *common.MostChanged {
	maxAddr := ""
	maxValue := big.NewInt(0)

	for addr, value := range addresses {
		absValue := value.Abs(value)
		if value.Cmp(maxValue) == 1 {
			maxAddr = addr
			maxValue = absValue
		}
	}

	return &common.MostChanged{
		Address: maxAddr,
		Value:   hexutil.EncodeBig(maxValue),
	}
}
