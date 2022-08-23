package main

import (
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/casper-ecosystem/casper-golang-sdk/keypair/ed25519"
	"github.com/casper-ecosystem/casper-golang-sdk/sdk"
	"github.com/casper-ecosystem/casper-golang-sdk/serialization"
	"github.com/casper-ecosystem/casper-golang-sdk/types"

	"test/helper"
)

func main() {
	rpcClient := sdk.NewRpcClient("http://localhost:11101/rpc")
	privKeyPath := "/home/jh/casper-node/utils/nctl/assets/net-1/users/user-1/secret_key.pem"
	pubKeyPath := "/home/jh/casper-node/utils/nctl/assets/net-1/users/user-1/public_key.pem"

	pair, _ := ed25519.ParseKeyFiles(pubKeyPath, privKeyPath)

	var hash32 [32]byte

	// contract hash
	// hash-c1f0f08d9a3cfc022e5baa5d6cbc645cd4a725503ab1df9eb5cc5b356788cbf8
	decodedHash, err2 := hex.DecodeString("c1f0f08d9a3cfc022e5baa5d6cbc645cd4a725503ab1df9eb5cc5b356788cbf8")
	if err2 != nil {
		return
	}

	for i := 0; i < 32; i++ {
		hash32[i] = decodedHash[i]
	}

	var argsOrder []string

	argsOrder = append(make([]string, 0), "CLTypeKeyAccount_value", "CLTypeKeyHash_value", "CLTypeKeyURef_value")
	// "CLTypeKeyURef_value")

	// KeyAccount
	// account-hash-004c61453f1bdf1f3c4b20b47b2fcfedabcc9e3afb29f8bb5983b7184e6a4497
	keyAccount := types.Key{
		Type:    types.KeyTypeAccount,
		Account: helper.GetAddress("004c61453f1bdf1f3c4b20b47b2fcfedabcc9e3afb29f8bb5983b7184e6a4497e5"),
	}
	KeyAccount_value, err := serialization.Marshal(types.CLValue{Type: types.CLTypeKey, Key: &keyAccount})

	// KeyTypeHash
	// hash-014c61453f1bdf1f3c4b20b47b2fcfedabcc9e3afb29f8bb5983b7184e6a4497
	keyHash := types.Key{
		Type: types.KeyTypeHash,
		Hash: helper.GetAddress("014c61453f1bdf1f3c4b20b47b2fcfedabcc9e3afb29f8bb5983b7184e6a4497e5"),
	}
	KeyHash_value, err := serialization.Marshal(types.CLValue{Type: types.CLTypeKey, Key: &keyHash})

	// URef_value
	// uref-024c61453f1bdf1f3c4b20b47b2fcfedabcc9e3afb29f8bb5983b7184e6a4497-007
	URef_value := types.URef{
		AccessRight: types.AccessRightReadAddWrite,
		Address:     helper.GetAddress("024c61453f1bdf1f3c4b20b47b2fcfedabcc9e3afb29f8bb5983b7184e6a4497e507"),
	}

	// KeyTypeURef
	keyURef := types.Key{
		Type: types.KeyTypeURef,
		URef: &URef_value}

	KeyURef_value, err := serialization.Marshal(types.CLValue{Type: types.CLTypeKey, Key: &keyURef})

	// ===============

	args := sdk.NewRunTimeArgs(map[string]sdk.Value{

		"CLTypeKeyAccount_value": {
			Tag:         types.CLTypeKey,
			StringBytes: hex.EncodeToString(KeyAccount_value),
		},
		"CLTypeKeyHash_value": {
			Tag:         types.CLTypeKey,
			StringBytes: hex.EncodeToString(KeyHash_value),
		},
		"CLTypeKeyURef_value": {
			Tag:         types.CLTypeKey,
			StringBytes: hex.EncodeToString(KeyURef_value),
		},
	}, argsOrder)

	session := sdk.NewStoredContractByHash(hash32, "transfer_token", *args)

	payment := sdk.StandardPayment(big.NewInt(10000000000))

	deploy := sdk.MakeDeploy(sdk.NewDeployParams(pair.PublicKey(), "casper-net-1", nil, 0), payment, session)

	deploy.SignDeploy(pair)
	result, err := rpcClient.PutDeploy(*deploy)
	if err != nil {
		fmt.Printf("err is:  %+v\n", err)
	}

	fmt.Printf("hash is => \n%+v\n", result.Hash)

}
