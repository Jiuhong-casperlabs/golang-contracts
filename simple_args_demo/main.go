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

	argsOrder = append(make([]string, 0), "CLTypeBool_value", "CLTypeI32_value", "CLTypeI64_value", "CLTypeU8_value", "CLTypeU32_value", "CLTypeU64_value", "CLTypeU128_value", "CLTypeU256_value", "CLTypeU512_value", "CLTypeUnit_value", "CLTypeString_value", "CLTypeURef_value", "CLTypePublicKey_value")

	var amount big.Int
	amount.SetUint64(600000000000)
	// ==== args =====

	// Bool_value
	Bool_value, _ := serialization.Marshal(true)

	// U512_value
	U512_value, _ := serialization.Marshal(serialization.U512{Int: amount})
	// U256_value
	U256_value, _ := serialization.Marshal(serialization.U256{Int: amount})
	// U128_value
	U128_value, _ := serialization.Marshal(serialization.U128{Int: amount})
	// U64_value
	U64_value, _ := serialization.Marshal(uint64(1024))
	// U32_value
	U32_value, _ := serialization.Marshal(uint32(1024))
	// U8_value
	U8_value, _ := serialization.Marshal(uint8(123))
	// int64(7)
	I64_value, _ := serialization.Marshal(int64(123))
	// int32(7)
	I32_value, _ := serialization.Marshal(int32(123))
	// unit_value
	unit_value, _ := serialization.Marshal(types.CLValue{Type: types.CLTypeUnit})
	// String_value
	String_value, _ := serialization.Marshal("helloworld")
	// URef_value
	// KeyTypeURef
	URef := types.URef{
		AccessRight: types.AccessRightReadAddWrite,
		Address:     helper.GetAddress("024c61453f1bdf1f3c4b20b47b2fcfedabcc9e3afb29f8bb5983b7184e6a4497e507"),
	}

	URef_value, _ := serialization.Marshal(types.CLValue{Type: types.CLTypeURef, URef: &URef})
	// public key
	publickey_value := "011542c5f1909889ac1f4937d9043c0f135fe229993f15780c45246a8d170617c7"
	// ===============

	args := sdk.NewRunTimeArgs(map[string]sdk.Value{
		"CLTypeBool_value": {
			Tag:         types.CLTypeBool,
			StringBytes: hex.EncodeToString(Bool_value),
		},
		"CLTypeI32_value": {
			Tag:         types.CLTypeI32,
			StringBytes: hex.EncodeToString(I32_value),
		},
		"CLTypeI64_value": {
			Tag:         types.CLTypeI64,
			StringBytes: hex.EncodeToString(I64_value),
		},
		"CLTypeU8_value": {
			Tag:         types.CLTypeU8,
			StringBytes: hex.EncodeToString(U8_value),
		},
		"CLTypeU32_value": {
			Tag:         types.CLTypeU32,
			StringBytes: hex.EncodeToString(U32_value),
		},
		"CLTypeU64_value": {
			Tag:         types.CLTypeU64,
			StringBytes: hex.EncodeToString(U64_value),
		},
		"CLTypeU128_value": {
			Tag:         types.CLTypeU128,
			StringBytes: hex.EncodeToString(U128_value),
		},
		"CLTypeU256_value": {
			Tag:         types.CLTypeU256,
			StringBytes: hex.EncodeToString(U256_value),
		},
		"CLTypeU512_value": {
			Tag:         types.CLTypeU512,
			StringBytes: hex.EncodeToString(U512_value),
		},
		"CLTypeUnit_value": {
			Tag:         types.CLTypeUnit,
			StringBytes: hex.EncodeToString(unit_value),
		},
		"CLTypeString_value": {
			Tag:         types.CLTypeString,
			StringBytes: hex.EncodeToString(String_value),
		},
		"CLTypeURef_value": {
			Tag:         types.CLTypeURef,
			StringBytes: hex.EncodeToString(URef_value),
		},
		"CLTypePublicKey_value": {
			Tag:         types.CLTypePublicKey,
			StringBytes: publickey_value,
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
