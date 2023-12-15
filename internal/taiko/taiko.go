package taiko

import (
	"taiko-tx/config"

	"github.com/ethereum/go-ethereum/common"
)

type Taiko struct {
	RPC                  string
	ApproveSmartContract common.Address
	DepositSmartContract common.Address
	Wallets              []config.Wallet
}

type Wallet struct {
	PrivateKey string
}

func Init(c *config.Configuration) *Taiko {
	return &Taiko{
		RPC:                  c.Taiko.RPC,
		ApproveSmartContract: common.HexToAddress(c.Taiko.ApproveSmartContract),
		DepositSmartContract: common.HexToAddress(c.Taiko.DepositSmartContract),
		Wallets:              c.Wallets,
	}
}
