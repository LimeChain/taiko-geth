package main

import (
	"os"
	"test-preconf/spammer"

	"github.com/charmbracelet/log"
)

func main() {
	logger := log.New(os.Stderr)
	logger.SetReportTimestamp(false)

	spammer := spammer.New(url, chainID, logger, accounts, maxTxsPerAccount)

	// send automatically generated txs
	// spammer.Start(txDefaults)

	// send manually preapred txs
	spammer.SendPreparedTxs(txsPerAccount, txDefaults)

	// sendSingleTx()
}
