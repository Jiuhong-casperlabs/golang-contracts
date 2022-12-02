package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	httpposturl := "http://3.136.227.9:7777/rpc"
	fmt.Println("HTTP JSON POST URL:", httpposturl)

	var jsonData = []byte(`{
		"id": 1,
		"jsonrpc": "2.0",
		"method": "chain_get_state_root_hash"
	}`)
	request, error := http.NewRequest("POST", httpposturl, bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, error := client.Do(request)
	if error != nil {
		panic(error)
	}
	defer response.Body.Close()

	fmt.Println("response Status:", response.Status)
	fmt.Println("response Headers:", response.Header)
	body, _ := ioutil.ReadAll(response.Body)
	fmt.Println("response Body:", string(body))

	fmt.Println("====")
	u := &State_root_hash_Result{}
	err := json.Unmarshal([]byte(string(body)), &u)
	// err := json.Unmarshal([]byte(encoded), &u)
	if err != nil {
		panic(err)
	}

	fmt.Printf("state_root_hash: %s\n", u.Result.State_root_hash)

}

type State_root_hash_Result struct {
    Jsonrpc string
	Id int
    Result struct {
        Api_version string
        State_root_hash string
    }
}

