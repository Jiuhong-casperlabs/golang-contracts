package main

import (
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/casper-ecosystem/casper-golang-sdk/keypair/ed25519"
	"github.com/casper-ecosystem/casper-golang-sdk/sdk"
	"github.com/casper-ecosystem/casper-golang-sdk/serialization"
	"github.com/casper-ecosystem/casper-golang-sdk/types"
)

func createPtrU32(v uint32) *uint32 {
	return &v
}

func createPtrList(v []types.CLValue) *[]types.CLValue {
	return &v
}

func toPtrU64(val uint64) *uint64 {
	return &val
}

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

	argsOrder = append(make([]string, 0), "list_value")
	//========== option_u64

	// List := types.List {}
	// URef := types.URef{
	// 	AccessRight: types.AccessRightReadAddWrite,
	// 	Address:     getAddress("024c61453f1bdf1f3c4b20b47b2fcfedabcc9e3afb29f8bb5983b7184e6a4497e507"),
	// }

	a := types.CLValue{Type: types.CLTypeList, List: createPtrList(make([]types.CLValue, 0))}
	list_bytes, err := serialization.Marshal(a)
	print("err...")
	print(err)
	list_value := sdk.Value{
		Tag:         types.CLTypeList,
		IsOptional:  false,
		StringBytes: hex.EncodeToString(list_bytes),
	}

	args := sdk.NewRunTimeArgs(map[string]sdk.Value{
		"list_value": list_value,
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
