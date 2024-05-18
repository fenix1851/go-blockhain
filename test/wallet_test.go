package wallet_test

import (
	"testing"

	wallet "github.com/fenix1851/golang-blockchain/wallet"
)

// TestWallet - это функция, которая тестирует функцию Wallet.
func TestNewWallet(t *testing.T) {
	wallet := wallet.NewWallet()
	address := wallet.Address()
	t.Logf("Address: %s", address)
}

// TestNewKeyPair - это функция, которая тестирует функцию NewKeyPair.

func TestNewKeyPair(t *testing.T) {
	priv, pub := wallet.NewKeyPair()
	if priv.X == nil {
		t.Error("Private key is nil")
		return
	}
	if pub == nil {
		t.Error("Public key is nil")
		return
	}
	t.Logf("Private key: %x", priv)
	t.Logf("Public key: %x", pub)
}

func TestCreateWallets(t *testing.T) {
	wallets, err := wallet.CreateWallets()
	if err != nil {
		t.Error(err)
	}
	t.Logf("Wallets: %v", wallets)
}

func TestWalletsAddWallet(t *testing.T) {
	wallets, err := wallet.CreateWallets()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	address := wallets.AddWallet()
	t.Logf("Address: %s", address)
}

func TestWalletsSaveFile(t *testing.T) {
	wallets, err := wallet.CreateWallets()
	if err != nil {
		t.Error(err)
	}
	err = wallets.SaveFile()
	if err != nil {
		t.Error(err)
	}
}

func TestWalletsGetAllAddresses(t *testing.T) {
	wallets, err := wallet.CreateWallets()
	if err != nil {
		t.Error(err)
	}
	addresses := wallets.GetAllAddresses()
	t.Logf("Addresses: %v", addresses)
}
