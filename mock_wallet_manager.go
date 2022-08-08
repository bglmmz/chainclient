package chainclient

import (
	"crypto/ecdsa"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"sync"
)

var (
	mockWalletOnce sync.Once
)

var mockWallet *MockWallet

type MockWallet struct {
	priKey        *ecdsa.PrivateKey
	pubKey        *ecdsa.PublicKey
	walletAddress common.Address
}

func MockWalletInstance() *MockWallet {
	return mockWallet
}

func InitMockWallet() {
	mockWalletOnce.Do(func() {
		mockWallet = new(MockWallet)

		key, _ := crypto.GenerateKey()
		mockWallet.priKey = key
		mockWallet.pubKey = &key.PublicKey
		mockWallet.walletAddress = crypto.PubkeyToAddress(key.PublicKey)
	})
}

// GetAddress returns the organization wallet address
func (m *MockWallet) GetAddress() common.Address {
	return m.walletAddress
}

// GetPrivateKey returns the organization private key
func (m *MockWallet) SetPrivateKey(privateKey *ecdsa.PrivateKey) {
	m.priKey = privateKey
	m.pubKey = &privateKey.PublicKey
	m.walletAddress = crypto.PubkeyToAddress(privateKey.PublicKey)
}

// GetPrivateKey returns the organization private key
func (m *MockWallet) GetPrivateKey() *ecdsa.PrivateKey {
	return m.priKey
}

// GetPrivateKey returns the organization private key
func (m *MockWallet) GetPublicKey() *ecdsa.PublicKey {
	return m.pubKey
}
