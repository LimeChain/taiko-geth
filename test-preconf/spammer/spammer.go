package spammer

import (
	"context"
	"math/big"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/charmbracelet/log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type Spammer struct {
	ctx              context.Context
	cancel           context.CancelFunc
	wg               *sync.WaitGroup
	client           *EthClient
	logger           *log.Logger
	accounts         []*Account
	maxTxsPerAccount uint64
}

func New(url string, chainID *big.Int, logger *log.Logger, accounts []*Account, maxTxsPerAccount uint64) *Spammer {
	var wg = new(sync.WaitGroup)

	ctx, cancel := context.WithCancel(context.Background())

	client, err := NewEthClient(ctx, url, chainID, logger)
	if err != nil {
		logger.Error("Failed to connect to the Ethereum client", "error", err)
	}

	return &Spammer{
		ctx:              ctx,
		cancel:           cancel,
		wg:               wg,
		client:           client,
		accounts:         accounts,
		logger:           logger,
		maxTxsPerAccount: maxTxsPerAccount,
	}
}

func (s *Spammer) Start(txDefaults func(i interface{}) *types.Transaction) {
	for {
		select {
		case <-s.ctx.Done():
			s.logger.Warn("Stopping spammer")
			return
		default:
			s.logger.Info("Sending new batch of txs\n")

			for _, account := range s.accounts {
				nonce, err := s.client.GetNonce(account)
				if err != nil {
					s.logger.Error("Failed to get nonce", "error", err)
					continue
				}
				account.Mutex.Lock()
				account.Nonce = nonce
				account.Mutex.Unlock()
			}

			assignedSlots, err := s.client.FetchAssignedSlots()
			if err != nil {
				s.logger.Error("Failed to fetch assigned slots", "error", err)
				continue
			}

			currentSlot, currentSlotEndTime, err := s.client.FetchCurrentSlot(time.Now().Unix())
			if err != nil {
				s.logger.Error("Failed to fetch current slot", "error", err)
				continue
			}

			firstAcceptableSlot := CalculateFirstAcceptableSlot(currentSlot, assignedSlots)
			deadline := new(big.Int).Add(firstAcceptableSlot, big.NewInt(0))

			s.sendTxs(txDefaults, deadline)

			durationUntilNextSlot := time.Duration(currentSlotEndTime-uint64(time.Now().Unix())) * time.Second
			s.logger.Info("Waiting until next slot", "duration", durationUntilNextSlot)
			time.Sleep(durationUntilNextSlot)
		}
	}
}

func (s *Spammer) sendTxs(assignTxDefaults func(i interface{}) *types.Transaction, deadline *big.Int) {
	for _, account := range s.accounts {
		s.wg.Add(1)

		go func(account *Account, deadline *big.Int) {
			defer s.wg.Done()

			for i := uint64(0); i < s.maxTxsPerAccount; i++ {
				tx := assignTxDefaults(types.InclusionPreconfirmationTx{
					Nonce:    account.Nonce,
					To:       account.Address(),
					Deadline: deadline,
				})

				// Sign and send the tx
				signedTx, err := s.client.SendTx(account, tx)
				if err != nil {
					s.logger.Error("Failed to send tx", "error", err)
					continue
				}
				account.Mutex.Lock()
				account.Nonce++
				account.Mutex.Unlock()

				// Log the tx
				s.client.LogTx(signedTx)

				// Log the receipt
				time.Sleep(100 * time.Millisecond)
				s.client.LogReceipt(signedTx)
			}
			s.logger.Info("All txs for account have been sent", "account", account.Address())

		}(account, deadline)
	}

	s.wg.Wait()
}

func (s *Spammer) SendPreparedTxs(
	txsPerAccount func(account *Account, nonce uint64, currentSlot uint64, assignedSlots []uint64) []interface{},
	txDefaults func(i interface{}) *types.Transaction) {
	for _, account := range s.accounts {
		func() {
			nonce, err := s.client.GetNonce(account)
			if err != nil {
				s.logger.Error("Failed to get nonce", "error", err)
			}

			currentSlot, _, err := s.client.FetchCurrentSlot(time.Now().Unix())
			if err != nil {
				s.logger.Error("Failed to fetch current slot", "error", err)
			}

			assignedSlots, err := s.client.FetchAssignedSlots()
			if err != nil {
				s.logger.Error("Failed to fetch assigned slots", "error", err)
			}

			// Get txs for the account
			txs := txsPerAccount(account, nonce, currentSlot, assignedSlots)

			// Iterate over each tx
			for _, t := range txs {
				tx := txDefaults(t)

				// Sign and send the tx
				signedTx, err := s.client.SendTx(account, tx)
				if err != nil {
					s.logger.Error("Failed to send tx", "error", err)
				}
				if signedTx != nil {
					// Log the tx
					s.client.LogTx(signedTx)

					// Log the receipt
					s.client.LogReceipt(signedTx)
				}
			}
		}()

		s.logger.Info("Done")
	}
}

func CalculateFirstAcceptableSlot(currentSlot uint64, assignedSlots []uint64) *big.Int {
	sort.Slice(assignedSlots, func(i, j int) bool {
		return assignedSlots[i] < assignedSlots[j]
	})

	logger := log.New(os.Stderr)
	logger.SetReportTimestamp(false)

	// logger.Error("Assigned slots", "slots", assignedSlots)
	// logger.Error("Current slot", "slot", currentSlot)

	var firstAcceptableSlot uint64
	for _, slot := range assignedSlots {
		if slot >= currentSlot+common.SlotsOffsetInAdvance {
			firstAcceptableSlot = slot
			break
		}
	}
	if firstAcceptableSlot == 0 {
		logger.Error("No acceptable slot found")
	}

	// logger.Info("Acceptable slot", "slot", firstAcceptableSlot)

	return new(big.Int).SetUint64(firstAcceptableSlot)
}
