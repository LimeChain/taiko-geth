package miner

import (
	"bytes"
	"compress/zlib"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/beacon/engine"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/txpool"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/holiman/uint256"
)

// sealBlockWith mines and seals a block with the given block metadata.
func (w *worker) sealBlockWith(
	parent common.Hash,
	timestamp uint64,
	blkMeta *engine.BlockMetadata,
	baseFeePerGas *big.Int,
	withdrawals types.Withdrawals,
) (*types.Block, error) {
	// Decode transactions bytes.
	var txs types.Transactions
	if err := rlp.DecodeBytes(blkMeta.TxList, &txs); err != nil {
		return nil, fmt.Errorf("failed to decode txList: %w", err)
	}

	if len(txs) == 0 {
		// A L2 block needs to have have at least one `V1TaikoL2.anchor` or
		// `V1TaikoL2.invalidateBlock` transaction.
		return nil, fmt.Errorf("too less transactions in the block")
	}

	params := &generateParams{
		timestamp:     timestamp,
		forceTime:     true,
		parentHash:    parent,
		coinbase:      blkMeta.Beneficiary,
		random:        blkMeta.MixHash,
		withdrawals:   withdrawals,
		noTxs:         false,
		baseFeePerGas: baseFeePerGas,
	}

	// Set extraData
	w.extra = blkMeta.ExtraData

	env, err := w.prepareWork(params)
	if err != nil {
		return nil, err
	}
	defer env.discard()

	env.header.GasLimit = blkMeta.GasLimit

	// Commit transactions.
	gasLimit := env.header.GasLimit
	rules := w.chain.Config().Rules(env.header.Number, true, timestamp)

	env.gasPool = new(core.GasPool).AddGas(gasLimit)

	for i, tx := range txs {
		// log.Warn("Seal tx", "index", i, "hash", tx.Hash().String())
		if i == 0 {
			if err := tx.MarkAsAnchor(); err != nil {
				return nil, err
			}
		}
		sender, err := types.LatestSignerForChainID(w.chainConfig.ChainID).Sender(tx)
		if err != nil {
			log.Error("Skip an invalid proposed transaction", "hash", tx.Hash(), "reason", err)
			continue
		}

		env.state.Prepare(rules, sender, blkMeta.Beneficiary, tx.To(), vm.ActivePrecompiles(rules), tx.AccessList())
		env.state.SetTxContext(tx.Hash(), env.tcount)
		if _, err := w.commitTransaction(env, tx); err != nil {
			log.Error("Skip an invalid proposed transaction", "hash", tx.Hash(), "reason", err)
			continue
		}
		env.tcount++
	}

	block, err := w.engine.FinalizeAndAssemble(w.chain, env.header, env.state, env.txs, nil, env.receipts, withdrawals)
	if err != nil {
		return nil, err
	}
	log.Warn("Seal block", "number", block.Number())

	results := make(chan *types.Block, 1)
	if err := w.engine.Seal(w.chain, block, results, nil); err != nil {
		return nil, err
	}
	block = <-results

	return block, nil
}

// CHANGE(limechain):

// 1. All transactions should all be able to pay the given base fee.
// 2. The total gas used should not exceed the given blockMaxGasLimit
// 3. The total bytes used should not exceed the given maxBytesPerTxList
// 4. The total number of transactions lists should not exceed the given maxTransactionsLists

// TODO(limechain): DRY UpdateTxSlotSnapshot and UpdateTxPoolSnapshot

func (w *worker) UpdateTxSlotSnapshot(
	snapshotSlot uint64,
	beneficiary common.Address,
	baseFee *big.Int,
	blockMaxGasLimit uint64,
	maxBytesPerTxList uint64,
	localAccounts []string,
	maxTransactionsLists uint64,
) error {
	currentHead := w.chain.CurrentBlock()
	if currentHead == nil {
		return fmt.Errorf("failed to find current head")
	}

	// Check if tx pool is empty at first.
	if len(w.getPendingTxs(baseFee)) == 0 {
		return nil
	}

	params := &generateParams{
		timestamp:     uint64(time.Now().Unix()),
		forceTime:     true,
		parentHash:    currentHead.Hash(),
		coinbase:      beneficiary,
		random:        currentHead.MixDigest,
		noTxs:         false,
		baseFeePerGas: baseFee,
	}

	env, err := w.prepareWork(params)
	if err != nil {
		return err
	}
	defer env.discard()

	var (
		signer = types.MakeSigner(w.chainConfig, new(big.Int).Add(currentHead.Number, common.Big1), currentHead.Time)

		// Fetch all preconf txs up to the current slot for simulation.

		// Split the pending transactions into locals and remotes, then
		// fill the block with all available pending transactions.
		localPreconfTxs, remotePreconfTxs = w.getPendingPreconfTxs(localAccounts, baseFee, snapshotSlot)
	)

	commitTxs := func() (*types.TxSlotSnapshot, error) {
		env.tcount = 0
		env.txs = []*types.Transaction{}
		env.gasPool = new(core.GasPool).AddGas(blockMaxGasLimit)
		env.header.GasLimit = blockMaxGasLimit

		var (
			localsPreconf  = make(map[common.Address][]*txpool.LazyTransaction)
			remotesPreconf = make(map[common.Address][]*txpool.LazyTransaction)
		)

		for address, txs := range localPreconfTxs {
			localsPreconf[address] = txs
		}
		for address, txs := range remotePreconfTxs {
			remotesPreconf[address] = txs
		}

		w.commitL2Transactions(
			env,
			newTransactionsByTypePriceAndNonce(signer, localsPreconf, baseFee),
			newTransactionsByTypePriceAndNonce(signer, remotesPreconf, baseFee),
			maxBytesPerTxList,
		)

		b, err := encodeAndComporeessTxList(env.txs)
		if err != nil {
			log.Error("Failed to encode and compress tx list", "err", err)
			return nil, err
		}

		totalBytes := uint64(len(b))
		totalGas := env.header.GasLimit - env.gasPool.Gas()

		// Calculate gas and bytes usage based on the difference
		// between current and previous slot snapshots.
		var (
			prevTotalGas, prevTotalBytes uint64
		)
		for i := uint64(0); i < snapshotSlot; i++ {
			txSlotSnapshot := w.txSnapshotsBuilder.GetTxSlotSnapshot(i)
			prevTotalGas += txSlotSnapshot.GasUsed
			prevTotalBytes += txSlotSnapshot.BytesLength
		}

		currentSlotGas := totalGas - prevTotalGas
		currentSlotBytes := totalBytes - prevTotalBytes

		if totalGas < prevTotalGas {
			currentSlotGas = 0
		}
		if totalBytes < prevTotalBytes {
			currentSlotBytes = 0
		}

		// all txs for specific slot have been executed, reset tx slot snapshot
		if currentSlotGas == 0 {
			w.txSnapshotsBuilder.ResetTxSlotSnapshot(snapshotSlot)
		}

		// Update the slot snapshot with only txs for the current slot.
		txSlotSnapshot := w.txSnapshotsBuilder.UpdateTxSlotSnapshot(snapshotSlot, env.txs, currentSlotBytes, currentSlotGas)

		// TODO(limechain): remove this log
		if txSlotSnapshot != nil && len(txSlotSnapshot.Txs) > 0 {
			log.Warn("Txs from slot snapshot", "slot", snapshotSlot, "tx count", len(txSlotSnapshot.Txs), "txs", txSlotSnapshot.Txs, "gas used", txSlotSnapshot.GasUsed, "bytes length", txSlotSnapshot.BytesLength)
		}

		return txSlotSnapshot, nil
	}

	for i := 0; i < int(maxTransactionsLists); i++ {
		txSlotSnapshot, err := commitTxs()
		if err != nil {
			log.Error("Failed to commit transactions", "err", err)
			return err
		}
		if txSlotSnapshot != nil && len(txSlotSnapshot.Txs) == 0 {
			break
		}
	}

	return nil
}

func (w *worker) UpdateTxPoolSnapshot(
	beneficiary common.Address,
	baseFee *big.Int,
	blockMaxGasLimit uint64,
	maxBytesPerTxList uint64,
	localAccounts []string,
	maxTransactionsLists uint64,
) error {
	currentHead := w.chain.CurrentBlock()
	if currentHead == nil {
		return fmt.Errorf("failed to find current head")
	}

	// Check if tx pool is empty at first.
	if len(w.getPendingTxs(baseFee)) == 0 {
		return nil
	}

	params := &generateParams{
		timestamp:     uint64(time.Now().Unix()),
		forceTime:     true,
		parentHash:    currentHead.Hash(),
		coinbase:      beneficiary,
		random:        currentHead.MixDigest,
		noTxs:         false,
		baseFeePerGas: baseFee,
	}

	env, err := w.prepareWork(params)
	if err != nil {
		return err
	}
	defer env.discard()

	var (
		signer = types.MakeSigner(w.chainConfig, new(big.Int).Add(currentHead.Number, common.Big1), currentHead.Time)
		// Split the pending transactions into locals and remotes, then
		// fill the block with all available pending transactions.
		localNonPreconfTxs, remoteNonPreconfTxs = w.getPendingNonPreconfTxs(localAccounts, baseFee)
	)

	commitTxs := func() (*types.TxPoolSnapshot, error) {
		env.tcount = 0
		env.txs = []*types.Transaction{}
		env.gasPool = new(core.GasPool).AddGas(blockMaxGasLimit)
		env.header.GasLimit = blockMaxGasLimit

		var (
			localsPreconf  = make(map[common.Address][]*txpool.LazyTransaction)
			remotesPreconf = make(map[common.Address][]*txpool.LazyTransaction)
		)

		for address, txs := range localNonPreconfTxs {
			localsPreconf[address] = txs
		}
		for address, txs := range remoteNonPreconfTxs {
			remotesPreconf[address] = txs
		}

		w.commitL2Transactions(
			env,
			newTransactionsByTypePriceAndNonce(signer, localsPreconf, baseFee),
			newTransactionsByTypePriceAndNonce(signer, remotesPreconf, baseFee),
			maxBytesPerTxList,
		)

		b, err := encodeAndComporeessTxList(env.txs)
		if err != nil {
			log.Error("Failed to encode and compress tx list", "err", err)
			return nil, err
		}

		gasUsed := env.header.GasLimit - env.gasPool.Gas()

		txPoolSnapshot := w.txSnapshotsBuilder.UpdateTxPoolSnapshot(env.txs, uint64(len(b)), gasUsed)

		// TODO(limechain): remove this log
		if txPoolSnapshot != nil && (len(txPoolSnapshot.PendingTxs) != 0 || len(txPoolSnapshot.ProposedTxs) != 0 || len(txPoolSnapshot.NewTxs) != 0) {
			log.Warn("Txs from pool snapshot", "pending tx count", len(txPoolSnapshot.PendingTxs), "txs", txPoolSnapshot.PendingTxs)
			log.Warn("Txs from pool snapshot", "proposed tx count", len(txPoolSnapshot.ProposedTxs), "txs", txPoolSnapshot.ProposedTxs)
			log.Warn("Txs from pool snapshot", "new tx count", len(txPoolSnapshot.NewTxs), "txs", txPoolSnapshot.NewTxs)
			log.Warn("Txs from pool snapshot", "gas used", txPoolSnapshot.GasUsed, "bytes length", txPoolSnapshot.BytesLength)
		}

		return txPoolSnapshot, nil
	}

	for i := 0; i < int(maxTransactionsLists); i++ {
		txPoolSnapshot, err := commitTxs()
		if err != nil {
			log.Error("Failed to commit transactions", "err", err)
			return err
		}

		if txPoolSnapshot != nil && len(txPoolSnapshot.NewTxs) == 0 {
			break
		}
	}

	return nil
}

func (w *worker) getPendingPreconfTxs(localAccounts []string, baseFee *big.Int, slotIndex uint64) (
	map[common.Address][]*txpool.LazyTransaction,
	map[common.Address][]*txpool.LazyTransaction,
) {
	pending := w.getPendingTxs(baseFee)

	preconfTxs := make(map[common.Address][]*txpool.LazyTransaction)

	// TODO(limechain): pass it from outside
	l1GenesisTimestamp := rawdb.ReadL1GenesisTimestamp(w.eth.BlockChain().DB())
	if l1GenesisTimestamp == nil {
		log.Error("L1 genesis timestamp is unknown")
	}

	for addr, txs := range pending {
		_, currentEpoch := common.CurrentSlotAndEpoch(*l1GenesisTimestamp, time.Now().Unix())
		for _, tx := range txs {
			txDeadlineSlot := tx.Tx.Deadline().Uint64()
			txDeadlineEpoch := txDeadlineSlot / uint64(common.EpochLength)

			if tx.Tx.Type() == types.InclusionPreconfirmationTxType && common.SlotIndex(txDeadlineSlot) <= slotIndex {
				if txDeadlineEpoch > currentEpoch {
					// do not pick this tx until the next epoch since we maintain only 32 slot snapshot for the current epoch
				} else {
					preconfTxs[addr] = append(preconfTxs[addr], tx)
				}
			}
		}
	}

	localPreconfTxs, remotePreconfTxs := make(map[common.Address][]*txpool.LazyTransaction), preconfTxs

	for _, local := range localAccounts {
		account := common.HexToAddress(local)
		if txs := remotePreconfTxs[account]; len(txs) > 0 {
			delete(remotePreconfTxs, account)
			localPreconfTxs[account] = txs
		}
	}

	return localPreconfTxs, remotePreconfTxs
}

func (w *worker) getPendingNonPreconfTxs(localAccounts []string, baseFee *big.Int) (
	map[common.Address][]*txpool.LazyTransaction,
	map[common.Address][]*txpool.LazyTransaction,
) {
	pending := w.getPendingTxs(baseFee)

	nonPreconfTxs := make(map[common.Address][]*txpool.LazyTransaction)

	for addr, txs := range pending {
		hasPreconfTxs := false

		for _, tx := range txs {
			if tx.Tx.Type() == types.InclusionPreconfirmationTxType {
				hasPreconfTxs = true
			} else {
				// if address has preconf txs before the non-preconf tx, then
				// delay the non-preconf txs until preconfs are processed.
				if !hasPreconfTxs {
					nonPreconfTxs[addr] = append(nonPreconfTxs[addr], tx)
				}
			}
		}
	}

	localTxs, remoteTxs := make(map[common.Address][]*txpool.LazyTransaction), nonPreconfTxs

	for _, local := range localAccounts {
		account := common.HexToAddress(local)
		if txs := remoteTxs[account]; len(txs) > 0 {
			delete(remoteTxs, account)
			localTxs[account] = txs
		}
	}

	return localTxs, remoteTxs
}

func (w *worker) getPendingTxs(baseFee *big.Int) map[common.Address][]*txpool.LazyTransaction {
	return w.eth.TxPool().Pending(txpool.PendingFilter{
		BaseFee:      uint256.MustFromBig(baseFee),
		OnlyPlainTxs: true},
	)
}

// commitL2Transactions tries to commit the transactions into the given state.
func (w *worker) commitL2Transactions(
	env *environment,
	txsLocal *transactionsByTypePriceAndNonce,
	txsRemote *transactionsByTypePriceAndNonce,
	maxBytesPerTxList uint64,
) {
	var (
		txs     = txsLocal
		isLocal = true
	)

	for {
		// If we don't have enough gas for any further transactions then we're done.
		if env.gasPool.Gas() < params.TxGas {
			log.Trace("Not enough gas for further transactions", "have", env.gasPool, "want", params.TxGas)
			break
		}

		// Retrieve the next transaction and abort if all done.
		ltx, _ := txs.Peek()
		if ltx == nil {
			if isLocal {
				txs = txsRemote
				isLocal = false
				continue
			}
			break
		}
		tx := ltx.Resolve()
		if tx == nil {
			log.Trace("Ignoring evicted transaction")

			txs.Pop()
			continue
		}

		if os.Getenv("TAIKO_MIN_TIP") != "" {
			minTip, err := strconv.Atoi(os.Getenv("TAIKO_MIN_TIP"))
			if err != nil {
				log.Error("Failed to parse TAIKO_MIN_TIP", "err", err)
			} else {
				if tx.GasTipCapIntCmp(new(big.Int).SetUint64(uint64(minTip))) < 0 {
					log.Trace("Ignoring transaction with low tip", "hash", tx.Hash(), "tip", tx.GasTipCap(), "minTip", minTip)
					txs.Pop()
					continue
				}
			}
		}

		// Error may be ignored here. The error has already been checked
		// during transaction acceptance is the transaction pool.
		from, _ := types.Sender(env.signer, tx)

		// Check whether the tx is replay protected. If we're not in the EIP155 hf
		// phase, start ignoring the sender until we do.
		if tx.Protected() && !w.chainConfig.IsEIP155(env.header.Number) {
			log.Trace("Ignoring reply protected transaction", "hash", tx.Hash(), "eip155", w.chainConfig.EIP155Block)

			txs.Pop()
			continue
		}
		// Start executing the transaction
		env.state.SetTxContext(tx.Hash(), env.tcount)

		// log.Warn("Processing tx", "from", from, "type", tx.Type(), "hash", tx.Hash(), "nonce", tx.Nonce(), "gasPrice", tx.GasPrice(), "gasTipCap", tx.GasTipCap(), "gasFeeCap", tx.GasFeeCap(), "gas", tx.Gas())
		_, err := w.commitTransaction(env, tx)
		switch {
		case errors.Is(err, core.ErrNonceTooLow):
			// New head notification data race between the transaction pool and miner, shift
			log.Trace("Skipping transaction with low nonce", "hash", ltx.Hash, "sender", from, "nonce", tx.Nonce())
			txs.Shift()

		case errors.Is(err, nil):
			// Everything ok, collect the logs and shift in the next transaction from the same account
			env.tcount++
			txs.Shift()

		default:
			// Transaction is regarded as invalid, drop all consecutive transactions from
			// the same sender because of `nonce-too-high` clause.
			log.Trace("Transaction failed, account skipped", "hash", ltx.Hash, "err", err)
			txs.Pop()
		}

		// Encode and compress the txList, if the byte length is > maxBytesPerTxList, remove the latest tx and break.
		b, err := encodeAndComporeessTxList(append(env.txs, tx))
		if err != nil {
			log.Trace("Failed to rlp encode and compress the pending transaction %s: %w", tx.Hash(), err)
			txs.Pop()
			continue
		}
		if len(b) > int(maxBytesPerTxList) {
			env.txs = env.txs[0 : env.tcount-1]
			break
		}
	}
	// log.Info("Committed transactions", "count", env.tcount)
}

// encodeAndComporeessTxList encodes and compresses the given transactions list.
func encodeAndComporeessTxList(txs types.Transactions) ([]byte, error) {
	b, err := rlp.EncodeToBytes(txs)
	if err != nil {
		return nil, err
	}

	return compress(b)
}

// compress compresses the given txList bytes using zlib.
func compress(txListBytes []byte) ([]byte, error) {
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	defer w.Close()

	if _, err := w.Write(txListBytes); err != nil {
		return nil, err
	}

	if err := w.Flush(); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}
