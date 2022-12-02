package main

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"test/helper"

	"github.com/casper-ecosystem/casper-golang-sdk/keypair/ed25519"
	"github.com/casper-ecosystem/casper-golang-sdk/sdk"
	"github.com/casper-ecosystem/casper-golang-sdk/serialization"
	"github.com/casper-ecosystem/casper-golang-sdk/types"
)

func main() {
	// rpcClient := sdk.NewRpcClient("http://localhost:11101/rpc")
	rpcClient := sdk.NewRpcClient("http://94.130.10.55:7777/rpc")
	// privKeyPath := "/home/jh/casper-node/utils/nctl/assets/net-1/users/user-1/secret_key.pem"
	// pubKeyPath := "/home/jh/casper-node/utils/nctl/assets/net-1/users/user-1/public_key.pem"

	privKeyPath := "/home/jh/keys/test1/secret_key.pem"
	pubKeyPath := "/home/jh/keys/test1/public_key.pem"

	pair, _ := ed25519.ParseKeyFiles(pubKeyPath, privKeyPath)

	var hash32 [32]byte

	// contract hash
	// hash-4120116565bd608fae6a45078055f320a2f429f426c86797b072b4efd15b186a
	decodedHash, err2 := hex.DecodeString("4120116565bd608fae6a45078055f320a2f429f426c86797b072b4efd15b186a")
	if err2 != nil {
		return
	}

	for i := 0; i < 32; i++ {
		hash32[i] = decodedHash[i]
	}

	var argsOrder []string

	argsOrder = append(make([]string, 0), "address",)

	// ==== args =====
	// KeyAccount
	// public-key: 0125a6336791eba195c472a8b7dbcd256a6ecddf8863e586a3dfefe2581a5d672c
	// account-hash-2293223427d59ebb331ac2221c3fcd1b3656a5cb72be924a6cdc9d52cdb6db0f
	keyAccount := types.Key{
		Type:    types.KeyTypeAccount,
		Account: helper.GetAddress("2293223427d59ebb331ac2221c3fcd1b3656a5cb72be924a6cdc9d52cdb6db0f"),
	}
	KeyAccount_value, err := serialization.Marshal(types.CLValue{Type: types.CLTypeKey, Key: &keyAccount})

	args := sdk.NewRunTimeArgs(map[string]sdk.Value{
		"address": {
			Tag:         types.CLTypeKey,
			StringBytes: hex.EncodeToString(KeyAccount_value),
		},
	}, argsOrder)

	session := sdk.NewStoredContractByHash(hash32, "balance_of", *args)

	payment := sdk.StandardPayment(big.NewInt(10000000000))

	deploy := sdk.MakeDeploy(sdk.NewDeployParams(pair.PublicKey(), "casper-test", nil, 0), payment, session)

	deploy.SignDeploy(pair)
	result, err := rpcClient.PutDeploy(*deploy)
	if err != nil {
		fmt.Printf("err is:  %+v\n", err)
	}

	fmt.Printf("hash is => \n%+v\n", result.Hash)
}
