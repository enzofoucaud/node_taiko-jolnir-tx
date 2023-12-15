package taiko

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"taiko-tx/internal/utils"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

var ctx = context.Background()

func (t *Taiko) BondTTKO(idWallet int) error {
	err := t.approve(idWallet)
	if err != nil {
		return err
	}

	err = t.depositTaikoToken(idWallet)
	if err != nil {
		return err
	}

	fmt.Println("Approved and deposited")

	return nil
}

func (t *Taiko) approve(idWallet int) error {
	privateKey, address, err := utils.GenerateECDSAKeys(t.Wallets[idWallet].PrivateKey)
	if err != nil {
		return err
	}

	client, err := ethclient.Dial(t.RPC)
	if err != nil {
		log.Fatal(err)
	}

	nonce, err := client.PendingNonceAt(context.Background(), address)
	if err != nil {
		return err
	}

	amountEth := big.NewInt(0)
	gasLimit := uint64(3000000)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return err
	}

	ChainID, err := client.NetworkID(ctx)
	if err != nil {
		return err
	}

	transferFnSignature := []byte("approve(address,uint256)")
	hash := crypto.Keccak256Hash(transferFnSignature)
	methodID := hash[:4]
	paddedAddress := common.LeftPadBytes(t.DepositSmartContract.Bytes(), 32)
	paddedAmount := common.LeftPadBytes(big.NewInt(2500000000000000000).Bytes(), 32)

	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAddress...)
	data = append(data, paddedAmount...)

	transaction := types.NewTransaction(nonce, t.ApproveSmartContract, amountEth, gasLimit, gasPrice, data)
	signedTx, err := types.SignTx(transaction, types.NewEIP155Signer(ChainID), privateKey)
	if err != nil {
		return err
	}
	err = client.SendTransaction(ctx, signedTx)
	if err != nil {
		return err
	}

	fmt.Printf("transaction sent: %s\n", signedTx.Hash().Hex())
	receipe, err := bind.WaitMined(ctx, client, signedTx)
	if err != nil {
		return err
	}
	fmt.Printf("transaction mined: %s\n", receipe.TxHash.Hex())

	return nil
}

func (t *Taiko) depositTaikoToken(idWallet int) error {
	privateKey, address, err := utils.GenerateECDSAKeys(t.Wallets[idWallet].PrivateKey)
	if err != nil {
		return err
	}

	client, err := ethclient.Dial(t.RPC)
	if err != nil {
		log.Fatal(err)
	}

	nonce, err := client.PendingNonceAt(context.Background(), address)
	if err != nil {
		return err
	}

	amountEth := big.NewInt(0)
	gasLimit := uint64(3000000)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return err
	}

	ChainID, err := client.NetworkID(ctx)
	if err != nil {
		return err
	}

	transferFnSignature := []byte("depositTaikoToken(uint256)")
	hash := crypto.Keccak256Hash(transferFnSignature)
	methodID := hash[:4]
	paddedAmount := common.LeftPadBytes(big.NewInt(2500000000000000000).Bytes(), 32)

	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAmount...)

	transaction := types.NewTransaction(nonce, t.DepositSmartContract, amountEth, gasLimit, gasPrice, data)
	signedTx, err := types.SignTx(transaction, types.NewEIP155Signer(ChainID), privateKey)
	if err != nil {
		return err
	}
	err = client.SendTransaction(ctx, signedTx)
	if err != nil {
		return err
	}

	fmt.Printf("transaction sent: %s\n", signedTx.Hash().Hex())
	receipe, err := bind.WaitMined(ctx, client, signedTx)
	if err != nil {
		return err
	}
	fmt.Printf("transaction mined: %s\n", receipe.TxHash.Hex())

	return nil
}
