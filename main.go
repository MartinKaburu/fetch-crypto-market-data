package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type cryptoExchangeAPI string

var cryptoExchangeAPIMap = map[string]cryptoExchangeAPI{"binance": "https://api1.binance.com/api/v3/exchangeInfo", "ftx": "https://ftx.com/api/markets", "kraken": "https://api.kraken.com/0/public/AssetPairs"}

func makeAPICall(url cryptoExchangeAPI, provider string) interface{} {
	r, err := http.Get(string(url))
	if err != nil {
		return err
	}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}
	switch provider {
	case "kraken":
		var data struct {
			Result map[string]struct {
				BaseCurrency  string `json:"base"`
				QuoteCurrency string `json:"quote"`
			} `json:"result"`
		}
		err = json.Unmarshal(body, &data)
		if err != nil {
			panic(err.Error())
		}
		return data
	case "ftx":
		var data struct {
			Result []struct {
				BaseCurrency  string `json:"baseCurrency"`
				QuoteCurrency string `json:"quoteCurrency"`
			} `json:"result"`
		}
		err = json.Unmarshal(body, &data)
		if err != nil {
			panic(err.Error())
		}
		return data
	case "binance":
		var data struct {
			Result []struct {
				BaseCurrency  string `json:"baseAsset"`
				QuoteCurrency string `json:"quoteAsset"`
			} `json:"symbols"`
		}
		err = json.Unmarshal(body, &data)
		if err != nil {
			panic(err.Error())
		}
		return data
	}
	var data interface{}
	return data
}

func main() {
	for provider, url := range cryptoExchangeAPIMap {
		fmt.Println(">>>>>>>>>>>>>>>>>>", provider, makeAPICall(url, provider))
	}

}
