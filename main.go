package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type cryptoExchangeAPI string

var cryptoExchangeAPIMap = map[string]cryptoExchangeAPI{"binance": "https://api1.binance.com/api/v3/exchangeInfo", "ftx": "https://ftx.com/api/markets", "kraken": "https://api.kraken.com/0/public/AssetPairs"}

type Currency struct {
	BaseCurrency  string
	QuoteCurrency string
}

func makeAPICall(url cryptoExchangeAPI, provider string) []Currency {
	r, err := http.Get(string(url))
	if err != nil {
		panic(err.Error())
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
		slc := make([]Currency, 0, len(data.Result))
		for _, ob := range data.Result {
			slc = append(slc, Currency(ob))
		}
		return slc
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
		slc := make([]Currency, 0, len(data.Result))
		for _, ob := range data.Result {
			slc = append(slc, Currency(ob))
		}
		return slc
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
		slc := make([]Currency, 0, len(data.Result))
		for _, ob := range data.Result {
			slc = append(slc, Currency(ob))
		}
		return slc
	}
	var data []Currency
	return data
}

func generatePairs(data []Currency) string {
	var buffer bytes.Buffer
	for _, obj := range data {
		buffer.WriteString(fmt.Sprintf("%s/%s\n", obj.BaseCurrency, obj.QuoteCurrency))
	}
	return buffer.String()
}

func main() {
	f, err := os.OpenFile("text.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err.Error())
	}
	for provider, url := range cryptoExchangeAPIMap {
		pairs := generatePairs(makeAPICall(url, provider))
		if _, err = f.WriteString(pairs); err != nil {
			panic(err.Error())
		}
	}

}
