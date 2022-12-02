package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

type State_root_hash_Result struct {
    Jsonrpc string
	Id int
    Result struct {
        Api_version string
        State_root_hash string
    }
}


func main() {
	httpposturl := "http://3.136.227.9:7777/rpc"
	fmt.Println("HTTP JSON POST URL:", httpposturl)

	// get
	var jsonData = []byte(`{
		"id": 1,
		"jsonrpc": "2.0",
		"method": "state_get_item",
		"params": {
			"state_root_hash": "6d59fdda440a55745b9b4882a77df19c32bac31a0e1483a6139a7ef11513a83e",
			"key": "hash-4120116565bd608fae6a45078055f320a2f429f426c86797b072b4efd15b186a"
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

	fmt.Println("response Status:", response.Status)
	fmt.Println("response Headers:", response.Header)
	body, _ := ioutil.ReadAll(response.Body)
	fmt.Println("response Body:", string(body))

}