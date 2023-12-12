package main

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/go-co-op/gocron"
)

var (
	scheduler       = gocron.NewScheduler(time.UTC)
	ctx             = context.Background()
	rpc             = "https://ethereum-sepolia.publicnode.com"
	contractApprove = common.HexToAddress("0x75F94f04d2144cB6056CCd0CFF1771573d838974")
	contractDeposit = common.HexToAddress("0x95fF8D3CE9dcB7455BEB7845143bEA84Fe5C4F6f")
	etherscanUrl    = "https://api-sepolia.etherscan.io/"
)

func main() {
	scheduler.Every(1).Hours().Do(verifyNode, os.Args[1], os.Args[2]) // nolint
	scheduler.StartBlocking()
}

func verifyNode(privKey string, etherscanApi string) {
	privateKey, err := crypto.HexToECDSA(privKey)
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("Public Key Error")
	}
	sender := crypto.PubkeyToAddress(*publicKeyECDSA)

	// http get request
	resp, err := http.Get(etherscanUrl + "api?module=account&action=txlist&page=1&offset=25&sort=desc&address=" + sender.Hex() + "&tag=latest&apikey=" + etherscanApi)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	type Response struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Result  []struct {
			TimeStamp string `json:"timeStamp"`
			Hash      string `json:"hash"`
			To        string `json:"to"`
		} `json:"result"`
	}

	var response Response

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		log.Fatal(err)
	}

	for _, tx := range response.Result {
		if tx.To == strings.ToLower(contractDeposit.Hex()) {
			if tx.TimeStamp == "" {
				fmt.Println("No timestamp, do nothing")
				continue
			}

			timestamp, err := strconv.ParseInt(tx.TimeStamp, 10, 64)
			if err != nil {
				log.Fatal(err)
			}
			t := time.Unix(timestamp, 0)

			if time.Since(t).Hours() > 3 {
				err = approve(privKey)
				if err != nil {
					log.Fatal(err)
				}

				err = depositTaikoToken(privKey)
				if err != nil {
					log.Fatal(err)
				}

				fmt.Println("Approved and deposited")
				break
			} else {
				fmt.Println("Timestamp:", t, "is less than 3 hours, do nothing")
				break
			}
		}
	}
}

func approve(privKey string) error {
	privateKey, err := crypto.HexToECDSA(privKey)
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return fmt.Errorf("Public Key Error")
	}
	sender := crypto.PubkeyToAddress(*publicKeyECDSA)

	client, err := ethclient.Dial(rpc)
	if err != nil {
		log.Fatal(err)
	}

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
	paddedAddress := common.LeftPadBytes(contractDeposit.Bytes(), 32)
	paddedAmount := common.LeftPadBytes(big.NewInt(2500000000000000000).Bytes(), 32)

	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAddress...)
	data = append(data, paddedAmount...)

	transaction := types.NewTransaction(nonce, contractApprove, amountEth, gasLimit, gasPrice, data)
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

func depositTaikoToken(privKey string) error {
	privateKey, err := crypto.HexToECDSA(privKey)
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return fmt.Errorf("Public Key Error")
	}
	sender := crypto.PubkeyToAddress(*publicKeyECDSA)

	client, err := ethclient.Dial(rpc)
	if err != nil {
		log.Fatal(err)
	}

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
	paddedAmount := common.LeftPadBytes(big.NewInt(2500000000000000000).Bytes(), 32)

	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAmount...)

	transaction := types.NewTransaction(nonce, contractDeposit, amountEth, gasLimit, gasPrice, data)
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
