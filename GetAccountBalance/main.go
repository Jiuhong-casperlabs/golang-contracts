package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"

	"github.com/casper-ecosystem/casper-golang-sdk/sdk"
	"github.com/pkg/errors"
)

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

type RpcClient struct {
	endpoint string
}

type balanceResponse struct {
	BalanceValue string `json:"balance_value"`
}

func (c *RpcClient) rpcCall(method string, params interface{}) (RpcResponse, error) {
	body, err := json.Marshal(RpcRequest{
		Version: "2.0",
		Method:  method,
		Params:  params,
	})

	if err != nil {
		return RpcResponse{}, errors.Wrap(err, "failed to marshal json")
	}

	resp, err := http.Post(c.endpoint, "application/json", bytes.NewReader(body))
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

func (c *RpcClient) GetAccountBalance(stateRootHash, balanceUref string) (big.Int, error) {
	resp, err := c.rpcCall("state_get_balance", map[string]string{
		"state_root_hash": stateRootHash,
		"purse_uref":      balanceUref,
	})
	if err != nil {
		return big.Int{}, err
	}

	var result balanceResponse
	err = json.Unmarshal(resp.Result, &result)
	if err != nil {
		return big.Int{}, fmt.Errorf("failed to get result: %w", err)
	}

	balance := big.Int{}
	balance.SetString(result.BalanceValue, 10)
	return balance, nil
}
func main() {
    var client = sdk.NewRpcClient("http://3.136.227.9:7777/rpc")
    stateRootHash := "c0eb76e0c3c7a928a0cb43e82eb4fad683d9ad626bcd3b7835a466c0587b0fff"
	key := "account-hash-a9efd010c7cee2245b5bad77e70d9beb73c8776cbe4698b2d8fdf6c8433d5ba0"

	balanceUref := client.GetAccountMainPurseURef(key)

	res, err := client.GetAccountBalance(stateRootHash, balanceUref)

	if err != nil {
		fmt.Errorf("failed to get result: %w", err)
	}

    fmt.Println(res)
}