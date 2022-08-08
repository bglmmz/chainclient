package chainclient

import (
	"context"
	"crypto/ecdsa"
	"github.com/PlatONnetwork/PlatON-Go/accounts/abi/bind"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/ethclient"
	"math/big"
	"time"
)

type Context interface {
	SetPrivateKey(*ecdsa.PrivateKey)
	GetPrivateKey() *ecdsa.PrivateKey
	GetPublicKey() *ecdsa.PublicKey
	GetAddress() common.Address
	GetClient() *ethclient.Client
	BlockNumber(timeoutCtx context.Context) (uint64, error)
	PendingNonceAt(timeoutCtx context.Context) (uint64, error)
	SuggestGasPrice(timeoutCtx context.Context) (*big.Int, error)

	EstimateGas(timeoutCtx context.Context, to common.Address, data []byte) (uint64, error)
	CallContract(timeoutCtx context.Context, to common.Address, data []byte) ([]byte, error)

	BuildTxOpts(value, gasLimit uint64) (*bind.TransactOpts, error)
	WaitReceipt(timeoutCtx context.Context, txHash common.Hash, interval time.Duration) *types.Receipt
	GetLog(timeoutCtx context.Context, toAddr common.Address, blockNo *big.Int) []*types.Log
}
