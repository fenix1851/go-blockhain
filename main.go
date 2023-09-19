package main

import (
	"os"

	"github.com/fenix1851/golang-blockchain/cli"
	"github.com/fenix1851/golang-blockchain/wallet"
)

func main() {
	defer os.Exit(0)
	w := wallet.NewWallet()
	w.Address()
	cmd := cli.CommandLine{}
	cmd.Run()
}
