package main

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"math/big"

	"github.com/casper-ecosystem/casper-golang-sdk/keypair/ed25519"
	"github.com/casper-ecosystem/casper-golang-sdk/sdk"
	"github.com/casper-ecosystem/casper-golang-sdk/serialization"
	"github.com/casper-ecosystem/casper-golang-sdk/types"
)

func main() {
	// nodeRpc := "http://159.65.118.250:7777/rpc"
	nodeRpc := "http://localhost:11101/rpc"
	privKeyPath := "/home/jh/casper-node/utils/nctl/assets/net-1/users/user-1/secret_key.pem"
	pubKeyPath := "/home/jh/casper-node/utils/nctl/assets/net-1/users/user-1/public_key.pem"
	modulePath := "/home/jh/casper-node/target/wasm32-unknown-unknown/release/undelegate.wasm"

	rpcClient := sdk.NewRpcClient(nodeRpc)

	pair, _ := ed25519.ParseKeyFiles(pubKeyPath, privKeyPath)
	module, _ := ioutil.ReadFile(modulePath)

	var argsOrder []string
	delegator := "011542c5f1909889ac1f4937d9043c0f135fe229993f15780c45246a8d170617c7"

	// set amount
	var amount big.Int
	amount.SetUint64(600000000000)

	amountBytes, _ := serialization.Marshal(serialization.U512{Int: amount})

	// set validator
	validator := "01db67635afa6ca726b398da1524e67413bd252f29484c25b07a0274103605d682"

	// set args order
	argsOrder = append(make([]string, 0), "amount", "delegator", "validator")

	args := sdk.NewRunTimeArgs(map[string]sdk.Value{
		"amount": {
			Tag:         types.CLTypeU512,
			StringBytes: hex.EncodeToString(amountBytes),
		},
		"delegator": {
			Tag:         types.CLTypePublicKey,
			StringBytes: delegator,
		},
		"validator": {
			Tag:         types.CLTypePublicKey,
			StringBytes: validator,
		},
	}, argsOrder)

	payment := sdk.StandardPayment(big.NewInt(10000000000))
	session := sdk.NewModuleBytes(module, *args)

	deploy := sdk.MakeDeploy(sdk.NewDeployParams(pair.PublicKey(), "casper-net-1", nil, 0), payment, session)
	deploy.SignDeploy(pair)

	result, err := rpcClient.PutDeploy(*deploy)
	if err != nil {
		fmt.Printf("err is:  %+v\n", err)
	}

	fmt.Printf("hash is => \n%+v\n", result.Hash)
}
