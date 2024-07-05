package preconfs

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/bindings"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/modern-go/reflect2"
)

const (
	defaultTimeout = 1 * time.Minute
)

// EthClient is a wrapper for go-ethereum eth client with a timeout attached.
type EthClient struct {
	chainID *big.Int
	timeout time.Duration
	*rpc.Client
	*ethClient
}

type ethClient struct {
	*ethclient.Client
}

func NewEthClient(ctx context.Context, url string, timeout time.Duration) (*EthClient, error) {
	var timeoutVal = defaultTimeout
	if timeout != 0 {
		timeoutVal = timeout
	}

	ethRpc, err := rpc.DialContext(ctx, url)
	if err != nil {
		return nil, err
	}

	ethClient := &ethClient{ethclient.NewClient(ethRpc)}

	chainID, err := ethClient.ChainID(ctx)
	if err != nil {
		return nil, err
	}

	return &EthClient{
		chainID:   chainID,
		ethClient: ethClient,
		timeout:   timeoutVal,
	}, nil
}

func (c *EthClient) CallContract(
	ctx context.Context,
	msg ethereum.CallMsg,
	blockNumber *big.Int,
) ([]byte, error) {
	ctxWithTimeout, cancel := ctxWithTimeoutOrDefault(ctx, c.timeout)
	defer cancel()

	return c.ethClient.CallContract(ctxWithTimeout, msg, blockNumber)
}

// ctxWithTimeoutOrDefault sets a context timeout if the deadline has not passed or is not set,
// and otherwise returns the context as passed in. cancel func is always set to an empty function
// so is safe to defer the cancel.
func ctxWithTimeoutOrDefault(ctx context.Context, defaultTimeout time.Duration) (context.Context, context.CancelFunc) {
	if reflect2.IsNil(ctx) {
		return context.WithTimeout(context.Background(), defaultTimeout)
	}
	if _, ok := ctx.Deadline(); !ok {
		return context.WithTimeout(ctx, defaultTimeout)
	}

	return ctx, func() {}
}

type Client struct {
	L1      *EthClient
	L2      *EthClient
	TaikoL2 *bindings.TaikoL2Client
}
