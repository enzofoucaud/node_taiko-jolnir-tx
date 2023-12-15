package config

import (
	"errors"

	"github.com/spf13/viper"
)

// Configurations exported
type Configuration struct {
	EtherScan EtherScan
	Taiko     Taiko
	Wallets   []Wallet
}

type EtherScan struct {
	Url    string
	ApiKey string
}

type Taiko struct {
	RPC                  string
	ApproveSmartContract string
	DepositSmartContract string
}

type Wallet struct {
	PrivateKey string
}

func Config() (*Configuration, error) {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("config file not found")
		}
		return nil, err
	}

	var wallets []Wallet
	err := viper.UnmarshalKey("WALLETS", &wallets)
	if err != nil {
		return nil, err
	}

	c := &Configuration{
		EtherScan: EtherScan{
			Url:    viper.GetString("etherScanUrl"),
			ApiKey: viper.GetString("etherScanApiKey"),
		},
		Taiko: Taiko{
			RPC:                  viper.GetString("taikoRpc"),
			ApproveSmartContract: viper.GetString("taikoApproveSmartContract"),
			DepositSmartContract: viper.GetString("taikoDepositSmartContract"),
		},
		Wallets: wallets,
	}

	return c, nil
}
