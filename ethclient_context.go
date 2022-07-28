package chainclient

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
	"time"
)

type EthContext struct {
	chainUrl string
	client   *ethclient.Client
	chainID  *big.Int
	wallet   Wallet
}

func NewEthClientContext(chainUrl string, wallet Wallet) *EthContext {
	ctx := new(EthContext)
	ctx.chainUrl = chainUrl
	client, err := ethclient.Dial(chainUrl)
	if err != nil {
		log.Fatal(err)
	}
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	ctx.client = client
	ctx.chainID = chainID
	return ctx
}

func (ctx *EthContext) GetClient() *ethclient.Client {
	return ctx.client
}
func (ctx *EthContext) GetPrivateKey() *ecdsa.PrivateKey {
	return ctx.wallet.GetPrivateKey()
}
func (ctx *EthContext) GetAddress() common.Address {
	return ctx.wallet.GetAddress()
}
func (ctx *EthContext) PendingNonceAt(timeoutCtx context.Context) (uint64, error) {
	return ctx.client.PendingNonceAt(timeoutCtx, ctx.GetAddress())
}
func (ctx *EthContext) SuggestGasPrice(timeoutCtx context.Context) (*big.Int, error) {
	return ctx.client.SuggestGasPrice(timeoutCtx)
}

func (ctx *EthContext) BlockNumber(timeoutCtx context.Context) (uint64, error) {
	return ctx.client.BlockNumber(timeoutCtx)
}

// EstimateGas uses context's wallet address as the caller (from)
func (ctx *EthContext) EstimateGas(timeoutCtx context.Context, to common.Address, input []byte) (uint64, error) {
	msg := ethereum.CallMsg{From: ctx.GetAddress(), To: &to, Data: input, Gas: 0, GasPrice: big.NewInt(0)}
	estimatedGas, err := ctx.client.EstimateGas(timeoutCtx, msg)
	if err != nil {
		return 0, err
	}
	return estimatedGas, nil
}

func (ctx *EthContext) CallContract(timeoutCtx context.Context, to common.Address, input []byte) ([]byte, error) {
	msg := ethereum.CallMsg{From: ctx.GetAddress(), To: &to, Data: input, Gas: 0, GasPrice: big.NewInt(0)}
	res, err := ctx.client.CallContract(timeoutCtx, msg, nil)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (ctx *EthContext) BuildTxOpts(value, gasLimit uint64) (*bind.TransactOpts, error) {
	nonce, err := ctx.client.PendingNonceAt(context.Background(), ctx.GetAddress())
	if err != nil {
		return nil, err
	}

	gasPrice, err := ctx.client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}

	//txOpts := bind.NewKeyedTransactor(m.Config.privateKey)
	txOpts, err := bind.NewKeyedTransactorWithChainID(ctx.GetPrivateKey(), ctx.chainID)
	if err != nil {
		return nil, err
	}

	txOpts.Nonce = new(big.Int).SetUint64(nonce)
	txOpts.Value = new(big.Int).SetUint64(value) // in wei
	txOpts.GasLimit = gasLimit                   // in units
	txOpts.GasPrice = gasPrice
	return txOpts, nil
}

func (ctx *EthContext) WaitReceipt(timeoutCtx context.Context, txHash common.Hash, interval time.Duration) *types.Receipt {
	fetchReceipt := func(txHash common.Hash) (*types.Receipt, error) {
		receipt, err := ctx.client.TransactionReceipt(context.Background(), txHash)
		if nil != err {
			//including NotFound
			log.Printf("failed to get transaction receipt, txHash: %s", txHash.Hex())
			return nil, err
		} else {
			log.Printf("transaction receipt: %#v", receipt)
			return receipt, nil
		}
	}

	if interval < 0 { // do once only
		receipt, err := fetchReceipt(txHash)
		if nil != err {
			return nil
		}
		return receipt

	} else {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-timeoutCtx.Done():
				log.Printf("get transaction receipt timeout, txHash: %s", txHash.Hex())
				return nil
			case <-ticker.C:
				if receipt, err := fetchReceipt(txHash); nil == err {
					return receipt
				}
			}
		}
	}
}

func (ctx *EthContext) GetLog(timeoutCtx context.Context, toAddr common.Address, blockNo *big.Int) []*types.Log {
	block, err := ctx.client.BlockByNumber(timeoutCtx, blockNo)
	if err != nil {
		log.Printf("get block error, block: %d, error: %v", blockNo.Uint64(), err)
		return nil
	}
	if block == nil {
		log.Printf("block not found, block: %d", blockNo.Uint64())
		return nil
	}

	logs := make([]*types.Log, 0)

	for _, tx := range block.Transactions() {
		if bytes.Compare(tx.To().Bytes(), toAddr.Bytes()) != 0 {
			continue
		}
		receipt, err := ctx.client.TransactionReceipt(timeoutCtx, tx.Hash())
		if err != nil {
			log.Printf("get tx receipt error, block: %d, txHash: %s, error: %v", blockNo.Uint64(), tx.Hash(), err)
			return nil
		}
		logs = append(logs, receipt.Logs...)
	}

	return logs
}
