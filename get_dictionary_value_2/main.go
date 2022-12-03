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
	"github.com/pkg/errors"
)

func rpcCall(method string, params interface{}) (RpcResponse, error) {
	httpposturl := "http://3.136.227.9:7777/rpc"

	body, err := json.Marshal(RpcRequest{
		Version: "2.0",
		Method:  method,
		Params:  params,
	})

	if err != nil {
		return RpcResponse{}, errors.Wrap(err, "failed to marshal json")
	}

	resp, err := http.Post(httpposturl, "application/json", bytes.NewReader(body))
	if err != nil {
		return RpcResponse{}, fmt.Errorf("failed to make request: %w", err)
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return RpcResponse{}, fmt.Errorf("failed to get response body: %w", err)
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return RpcResponse{}, fmt.Errorf("request failed, status code - %d, response - %s", resp.StatusCode, string(b))
	}

	var rpcResponse RpcResponse
	err = json.Unmarshal(b, &rpcResponse)
	if err != nil {
		return RpcResponse{}, fmt.Errorf("failed to parse response body: %w", err)
	}

	if rpcResponse.Error != nil {
		return rpcResponse, fmt.Errorf("rpc call failed, code - %d, message - %s", rpcResponse.Error.Code, rpcResponse.Error.Message)
	}

	return rpcResponse, nil
}
func main() {

	// contract hash "hash-4120116565bd608fae6a45078055f320a2f429f426c86797b072b4efd15b186a"
	const contract_hash = "hash-4120116565bd608fae6a45078055f320a2f429f426c86797b072b4efd15b186a"
	// public_key: 0125a6336791eba195c472a8b7dbcd256a6ecddf8863e586a3dfefe2581a5d672c
	const publKey = "25a6336791eba195c472a8b7dbcd256a6ecddf8863e586a3dfefe2581a5d672c"
	publKeyBytes, err := hex.DecodeString(publKey)
	if err != nil {
		return
	}

	// account hash: account-hash-2293223427d59ebb331ac2221c3fcd1b3656a5cb72be924a6cdc9d52cdb6db0f
	accountHex := ed25519.AccountHash(publKeyBytes)

	resultHex, err := hex.DecodeString(accountHex)
	if err != nil {
		fmt.Print(err)
	}

	// === step1 get dictionary_item_key ===
	zero := []byte{0}
	item_key_bytes := append(zero, resultHex...)
	item_key := b64.StdEncoding.EncodeToString(item_key_bytes)

	// === step2 get state-root-hash===
	resp, err := rpcCall("chain_get_state_root_hash",nil)
	b, _ := json.Marshal(resp)
	if err != nil {
		fmt.Println(err)
	}
	var result1 State_root_hash_Result
	err = json.Unmarshal(b, &result1)
	if err != nil {
		fmt.Println(err)
	}
	state_root_hash := result1.Result.State_root_hash

	// === step3 get rpc ===========
	contractNamedKey := map[string]string{
		"key":                 contract_hash, //contract hash
		"dictionary_name":     "balances",                                                              //dictionary uref name
		"dictionary_item_key": item_key,                          // dictionary item key
	}
	dictionary_identifier := map[string]interface{}{
		"ContractNamedKey": contractNamedKey,
	}
	resp, err = rpcCall("state_get_dictionary_item", map[string]interface{}{
		"state_root_hash":       state_root_hash,
		"dictionary_identifier": dictionary_identifier,
	})

	b, _ = json.Marshal(resp)

	if err != nil {
		fmt.Println(err)
	}

	var result Dictionary_value_Result
	err = json.Unmarshal(b, &result)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("balance=> %s\n", result.Result.Stored_value.CLValue.Parsed)

}
type State_root_hash_Result struct {
    Jsonrpc string
	Id string
    Result struct {
        Api_version string
        State_root_hash string
    }
}

type Dictionary_value_Result struct {
	Jsonrpc string
	Id      string
	Result  struct {
		Api_version    string
		Dictionary_key string
		Stored_value   struct {
			CLValue struct {
				Cl_type string
				Bytes   string
				Parsed  string
			}
		}
	}
}
type RpcRequest struct {
	Version string      `json:"jsonrpc"`
	Id      string      `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
}

type RpcResponse struct {
	Version string          `json:"jsonrpc"`
	Id      string          `json:"id"`
	Result  json.RawMessage `json:"result"`
	Error   *RpcError       `json:"error,omitempty"`
}

type RpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}


