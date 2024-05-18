package wallet

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"os"
)

const walletFile = "./tmp/wallets.data"

type Wallets struct {
	Wallets map[string]*Wallet
}

func init() {
	gob.Register(elliptic.P256()) // Register the elliptic curve type
}

func (ws *Wallets) LoadFile() error {
	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		return err
	}

	fileContent, err := os.ReadFile(walletFile)
	if err != nil {
		return err
	}

	var walletWrappers map[string]*WalletWrapper
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&walletWrappers)
	if err != nil {
		return err
	}

	ws.Wallets = make(map[string]*Wallet)
	for address, wrapper := range walletWrappers {
		ws.Wallets[address] = wrapper.ToWallet()
	}

	return nil
}

func (ws *Wallets) GetAllAddresses() []string {
	var addresses []string
	for address := range ws.Wallets {
		addresses = append(addresses, string(address))
	}
	return addresses
}

func (ws *Wallets) GetWallet(address string) *Wallet {
	return ws.Wallets[address]
}

func CreateWallets() (*Wallets, error) {
	wallets := Wallets{}
	wallets.Wallets = make(map[string]*Wallet)
	err := wallets.LoadFile()
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	if err != nil {
		fmt.Println("No wallet file found. Creating new wallet file...")
	}
	if os.IsNotExist(err) {
		err = wallets.SaveFile()
		if err != nil {
			return nil, err
		}
	}

	return &wallets, nil
}

func (ws *Wallets) AddWallet() string {
	wallet := NewWallet()
	address := string(wallet.Address())
	ws.Wallets[address] = wallet
	err := ws.SaveFile()
	if err != nil {
		fmt.Println(err)
	}

	return address
}

func (ws *Wallets) SaveFile() error {
	var content bytes.Buffer
	walletWrappers := make(map[string]*WalletWrapper)
	for address, wallet := range ws.Wallets {
		walletWrappers[address] = NewWalletWrapper(wallet)
	}
	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(walletWrappers)
	if err != nil {
		return err
	}
	err = os.WriteFile(walletFile, content.Bytes(), 0644)
	return err
}
