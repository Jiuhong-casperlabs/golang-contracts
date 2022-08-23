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

	argsOrder = append(make([]string, 0), "option_bool", "option_u8", "Option_u64", "Option_U128", "option_string", "option_key", "option_uref")

	// ==========option_bool
	option_bool := sdk.Value{
		Tag:        types.CLTypeOption,
		IsOptional: true,
		Optional: &sdk.Value{
			Tag: types.CLTypeBool,
		},
	}

	id_bool := true
	idBytes_bool, err := serialization.Marshal(types.CLValue{Type: types.CLTypeOption, Option: &types.CLValue{Type: types.CLTypeBool, Bool: &id_bool}})

	option_bool.Optional.StringBytes = hex.EncodeToString(idBytes_bool)

	// ======option_U8
	option_u8 := sdk.Value{
		Tag:        types.CLTypeOption,
		IsOptional: true,
		Optional: &sdk.Value{
			Tag: types.CLTypeU8,
		},
	}
	id_u8 := uint8(4)
	idBytes_u8, err := serialization.Marshal(types.CLValue{Type: types.CLTypeOption, Option: &types.CLValue{Type: types.CLTypeU8, U8: &id_u8}})

	option_u8.Optional.StringBytes = hex.EncodeToString(idBytes_u8)

	//========== option_u64
	option_u64 := sdk.Value{
		Tag:        types.CLTypeOption,
		IsOptional: true,
		Optional: &sdk.Value{
			Tag: types.CLTypeU64,
		},
	}

	id_u64 := uint64(1)
	idBytes_u64, err := serialization.Marshal(types.CLValue{Type: types.CLTypeOption, Option: &types.CLValue{Type: types.CLTypeU64, U64: &id_u64}})

	option_u64.Optional.StringBytes = hex.EncodeToString(idBytes_u64)

	//========== option U128
	option_U128 := sdk.Value{
		Tag:        types.CLTypeOption,
		IsOptional: true,
		Optional: &sdk.Value{
			Tag: types.CLTypeU128,
		},
	}
	var amount big.Int
	amount.SetUint64(600000000000)

	idBytes_U128, err := serialization.Marshal(types.CLValue{Type: types.CLTypeOption, Option: &types.CLValue{Type: types.CLTypeU128, U128: &amount}})

	option_U128.Optional.StringBytes = hex.EncodeToString(idBytes_U128)

	//========== option string
	option_string := sdk.Value{
		Tag:        types.CLTypeOption,
		IsOptional: true,
		Optional: &sdk.Value{
			Tag: types.CLTypeString,
		},
	}

	var hello = "hello"
	idBytes_String, err := serialization.Marshal(types.CLValue{Type: types.CLTypeOption, Option: &types.CLValue{Type: types.CLTypeString, String: &hello}})

	option_string.Optional.StringBytes = hex.EncodeToString(idBytes_String)

	//========== option key
	option_key := sdk.Value{
		Tag:        types.CLTypeOption,
		IsOptional: true,
		Optional: &sdk.Value{
			Tag: types.CLTypeKey,
		},
	}
	keyHash := types.Key{
		Type:    types.KeyTypeAccount,
		Account: helper.GetAddress("004c61453f1bdf1f3c4b20b47b2fcfedabcc9e3afb29f8bb5983b7184e6a4497e5"),
	}
	KeyHash_value, err := serialization.Marshal(types.CLValue{Type: types.CLTypeOption, Option: &types.CLValue{Type: types.CLTypeKey, Key: &keyHash}})
	option_key.Optional.StringBytes = hex.EncodeToString(KeyHash_value)

	//========== option uref
	// uref-024c61453f1bdf1f3c4b20b47b2fcfedabcc9e3afb29f8bb5983b7184e6a4497-007
	option_uref := sdk.Value{
		Tag:        types.CLTypeOption,
		IsOptional: true,
		Optional: &sdk.Value{
			Tag: types.CLTypeURef,
		},
	}

	URef := types.URef{
		AccessRight: types.AccessRightReadAddWrite,
		Address:     helper.GetAddress("024c61453f1bdf1f3c4b20b47b2fcfedabcc9e3afb29f8bb5983b7184e6a4497e507"),
	}

	URef_value, _ := serialization.Marshal(types.CLValue{Type: types.CLTypeOption, Option: &types.CLValue{Type: types.CLTypeURef, URef: &URef}})

	option_uref.Optional.StringBytes = hex.EncodeToString(URef_value)
	// ===
	// args
	args := sdk.NewRunTimeArgs(map[string]sdk.Value{
		"option_bool":   option_bool,
		"option_u8":     option_u8,
		"Option_u64":    option_u64,
		"Option_U128":   option_U128,
		"option_string": option_string,
		"option_key":    option_key,
		"option_uref":   option_uref,
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
