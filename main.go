package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	ctx             = context.Background()
	rpc             = "https://ethereum-sepolia.publicnode.com"
	contractApprove = common.HexToAddress("0x75F94f04d2144cB6056CCd0CFF1771573d838974")
	contractDeposit = common.HexToAddress("0x95fF8D3CE9dcB7455BEB7845143bEA84Fe5C4F6f")
)

func main() {
	privKey := os.Args[1]

	privateKey, err := crypto.HexToECDSA(privKey)
	if err != nil {
		log.Fatal(err)
	}

	err = approve(privateKey)
	if err != nil {
		log.Fatal(err)
	}

	err = depositTaikoToken(privateKey)
	if err != nil {
		log.Fatal(err)
	}
}

func approve(privateKey *ecdsa.PrivateKey) error {
	client, err := ethclient.Dial(rpc)
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return fmt.Errorf("Public Key Error")
	}
	sender := crypto.PubkeyToAddress(*publicKeyECDSA)

	nonce, err := client.PendingNonceAt(context.Background(), sender)
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
	fmt.Println("methodID", hexutil.Encode(methodID))

	paddedAddress := common.LeftPadBytes(contractDeposit.Bytes(), 32)
	fmt.Println("paddedAddress", hexutil.Encode(paddedAddress))

	paddedAmount := common.LeftPadBytes(big.NewInt(2500000000000000000).Bytes(), 32)
	fmt.Println("paddedAmount", hexutil.Encode(paddedAmount))

	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAddress...)
	data = append(data, paddedAmount...)

	transaction := types.NewTransaction(nonce, contractApprove, amountEth, uint64(gasLimit), gasPrice, data)
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

func depositTaikoToken(privateKey *ecdsa.PrivateKey) error {
	client, err := ethclient.Dial(rpc)
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return fmt.Errorf("Public Key Error")
	}
	sender := crypto.PubkeyToAddress(*publicKeyECDSA)

	nonce, err := client.PendingNonceAt(context.Background(), sender)
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
	fmt.Println("methodID", hexutil.Encode(methodID))

	paddedAmount := common.LeftPadBytes(big.NewInt(2500000000000000000).Bytes(), 32)
	fmt.Println("paddedAmount", hexutil.Encode(paddedAmount))

	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAmount...)

	transaction := types.NewTransaction(nonce, contractDeposit, amountEth, uint64(gasLimit), gasPrice, data)
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
