package main

import (
	"bytes"
	b64 "encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/casper-ecosystem/casper-golang-sdk/keypair/ed25519"
)

func main() {
	// get dictionary_item_key
	// public_key: 01d995c93ac47e763433b5ec973cac464c7343d76d6bd47c936cf8ce5d83032061
	// account hash: account-hash-448d833d4c5883a1be55cc3db63afbf8ac320b6d506fe80c7221e9db1d5ff699
	const publKey = "25a6336791eba195c472a8b7dbcd256a6ecddf8863e586a3dfefe2581a5d672c"
	publKeyBytes, err := hex.DecodeString(publKey)
	if err != nil {
		return
	}
	// fmt.Print(AccountHex(publKeyBytes))
	accountHex := ed25519.AccountHash(publKeyBytes)
	fmt.Println(accountHex)
	
	resultHex, err := hex.DecodeString(accountHex)
	if err != nil {
		fmt.Print(err) 
	}

    zero := []byte{0}
    five := append(zero,resultHex...)
	fmt.Println(five)
	sEnc := b64.StdEncoding.EncodeToString(five)
    fmt.Println(sEnc)
	// get rpc 
	httpposturl := "http://3.136.227.9:7777/rpc"
	fmt.Println("HTTP JSON POST URL:", httpposturl)

	var jsonData = []byte(`{
		"id": 1,
		"jsonrpc":"2.0",
		"method":"state_get_dictionary_item",
		 "params": {
			"state_root_hash": "6d59fdda440a55745b9b4882a77df19c32bac31a0e1483a6139a7ef11513a83e",
			"dictionary_identifier": {
			  "ContractNamedKey": {
				"key": "hash-4120116565bd608fae6a45078055f320a2f429f426c86797b072b4efd15b186a",
				"dictionary_name": "balances",
				"dictionary_item_key": "ACKTIjQn1Z67MxrCIhw/zRs2VqXLcr6SSmzcnVLNttsP"
			  }
			}
		  }
	}`)
	request, error := http.NewRequest("POST", httpposturl, bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, error := client.Do(request)
	if error != nil {
		panic(error)
	}
	defer response.Body.Close()

	// fmt.Println("response Status:", response.Status)
	// fmt.Println("response Headers:", response.Header)
	body, _ := ioutil.ReadAll(response.Body)
	// fmt.Println("response Body:", string(body))

	// fmt.Println("====")
	u := &Dictionary_value_Result{}
	err = json.Unmarshal([]byte(string(body)), &u)
	// err := json.Unmarshal([]byte(encoded), &u)
	if err != nil {
		panic(err)
	}

	fmt.Printf("balance: %s\n", u.Result.Stored_value.CLValue.Parsed)

}

type Dictionary_value_Result struct {
    Jsonrpc string
	Id int
    Result struct {
        Api_version string
        Dictionary_key string
		Stored_value struct {
			CLValue struct {
				Cl_type string
				Bytes string
				Parsed string
			}
			// Merkle_proof string
		}
    }
}

