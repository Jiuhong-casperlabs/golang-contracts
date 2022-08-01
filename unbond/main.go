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
	modulePath := "/home/jh/casper-node/target/wasm32-unknown-unknown/release/withdraw_bid.wasm"

	rpcClient := sdk.NewRpcClient(nodeRpc)

	pair, _ := ed25519.ParseKeyFiles(pubKeyPath, privKeyPath)
	module, _ := ioutil.ReadFile(modulePath)

	// set public key
	public_key := "011542c5f1909889ac1f4937d9043c0f135fe229993f15780c45246a8d170617c7"

	// set amount
	amount := big.NewInt(10000000000)
	amountBytes, _ := serialization.Marshal(serialization.U512{Int: *amount})

	// set unbond_purse
	unbond_purse := sdk.Value{
		Tag:        types.CLTypeOption,
		IsOptional: true,
		Optional: &sdk.Value{
			Tag: types.CLTypeURef,
		},
	}

	idBytes, _ := serialization.Marshal(types.CLValue{Type: types.CLTypeOption})

	unbond_purse.Optional.StringBytes = hex.EncodeToString(idBytes)

	// set args order
	var argsOrder []string
	argsOrder = append(make([]string, 0), "amount", "public_key", "unbond_purse")

	// set args
	args := sdk.NewRunTimeArgs(map[string]sdk.Value{
		"amount": {
			Tag:         types.CLTypeU512,
			StringBytes: hex.EncodeToString(amountBytes),
		},
		"public_key": {
			Tag:         types.CLTypePublicKey,
			StringBytes: public_key,
		},
		"unbond_purse": unbond_purse,
	}, argsOrder)

	// set payment
	payment := sdk.StandardPayment(big.NewInt(10000000000))
	// set session
	session := sdk.NewModuleBytes(module, *args)

	deploy := sdk.MakeDeploy(sdk.NewDeployParams(pair.PublicKey(), "casper-net-1", nil, 0), payment, session)
	deploy.SignDeploy(pair)

	result, err := rpcClient.PutDeploy(*deploy)
	if err != nil {
		fmt.Printf("err is:  %+v\n", err)
	}

	fmt.Printf("hash is => \n%+v\n", result.Hash)
}
