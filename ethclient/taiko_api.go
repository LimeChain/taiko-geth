package ethclient

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/eth"
)

// HeadL1Origin returns the latest L2 block's corresponding L1 origin.
func (ec *Client) HeadL1Origin(ctx context.Context) (*rawdb.L1Origin, error) {
	var res *rawdb.L1Origin

	if err := ec.c.CallContext(ctx, &res, "taiko_headL1Origin"); err != nil {
		return nil, err
	}

	return res, nil
}

// L1OriginByID returns the L2 block's corresponding L1 origin.
func (ec *Client) L1OriginByID(ctx context.Context, blockID *big.Int) (*rawdb.L1Origin, error) {
	var res *rawdb.L1Origin

	if err := ec.c.CallContext(ctx, &res, "taiko_l1OriginByID", hexutil.EncodeBig(blockID)); err != nil {
		return nil, err
	}

	return res, nil
}

// GetSyncMode returns the current sync mode of the L2 node.
func (ec *Client) GetSyncMode(ctx context.Context) (string, error) {
	var res string

	if err := ec.c.CallContext(ctx, &res, "taiko_getSyncMode"); err != nil {
		return "", err
	}

	return res, nil
}

func (ec *Client) GetPreconfirmedVirtualBlock(ctx context.Context) (eth.HashAndNumber, error) {
	var res eth.HashAndNumber

	if err := ec.c.CallContext(ctx, &res, "taiko_getPreconfirmedVirtualBlock"); err != nil {
		return eth.HashAndNumber{}, err
	}

	return res, nil
}

func (ec *Client) GetPendingVirtualBlock(ctx context.Context) (eth.HashAndNumber, error) {
	var res eth.HashAndNumber

	if err := ec.c.CallContext(ctx, &res, "taiko_getPendingVirtualBlock"); err != nil {
		return eth.HashAndNumber{}, err
	}

	return res, nil
}

func (ec *Client) UpdatePreconfirmedVirtualBlock(ctx context.Context, hash common.Hash, number *big.Int) (bool, error) {
	var res bool

	if err := ec.c.CallContext(ctx, &res, "taiko_updatePreconfirmedVirtualBlock", hash.Hex(), hexutil.EncodeBig(number)); err != nil {
		return false, err
	}

	return res, nil
}

func (ec *Client) UpdatePendingVirtualBlock(ctx context.Context, hash common.Hash, number *big.Int) (bool, error) {
	var res bool

	if err := ec.c.CallContext(ctx, &res, "taiko_updatePendingVirtualBlock", hash.Hex(), hexutil.EncodeBig(number)); err != nil {
		return false, err
	}

	return res, nil
}

func (ec *Client) DeletePendingVirtualBlock(ctx context.Context) error {
	var res bool

	if err := ec.c.CallContext(ctx, &res, "taiko_deletePendingVirtualBlock"); err != nil {
		return err
	}

	return nil
}
