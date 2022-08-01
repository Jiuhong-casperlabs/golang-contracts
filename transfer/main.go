package main

import (
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/casper-ecosystem/casper-golang-sdk/keypair"
	"github.com/casper-ecosystem/casper-golang-sdk/keypair/ed25519"
	"github.com/casper-ecosystem/casper-golang-sdk/sdk"
)

func main() {
	// nodeRpc := "http://159.65.118.250:7777/rpc"
	nodeRpc := "http://localhost:11101/rpc"
	rpcClient := sdk.NewRpcClient(nodeRpc)

	// set source key path
	privKeyPath := "/home/jh/casper-node/utils/nctl/assets/net-1/users/user-1/secret_key.pem"
	pubKeyPath := "/home/jh/casper-node/utils/nctl/assets/net-1/users/user-1/public_key.pem"

	pair, _ := ed25519.ParseKeyFiles(pubKeyPath, privKeyPath)

	// set target public key "01272a2fe949347aa893fdcbb99bfeb4c57e348c5359a45363514c4e15364e5136"
	decodedDest, _ := hex.DecodeString("272a2fe949347aa893fdcbb99bfeb4c57e348c5359a45363514c4e15364e5136")

	dest := &keypair.PublicKey{
		Tag:        keypair.KeyTagEd25519,
		PubKeyData: decodedDest,
	}

	// set payment
	payment := sdk.StandardPayment(big.NewInt(10000000000))
	// set session
	session := sdk.NewTransfer(big.NewInt(2500000000), dest, "", uint64(1))

	deploy := sdk.MakeDeploy(sdk.NewDeployParams(pair.PublicKey(), "casper-net-1", nil, 0), payment, session)
	deploy.SignDeploy(pair)

	result, err := rpcClient.PutDeploy(*deploy)
	if err != nil {
		fmt.Printf("err is:  %+v\n", err)
	}

	fmt.Printf("hash is => \n%+v\n", result.Hash)
}
