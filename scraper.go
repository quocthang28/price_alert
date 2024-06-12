package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

var errStr string = "error while fetching prices: %s"

func getCryptoPrices(symbols []string) string {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://pro-api.coinmarketcap.com/v2/cryptocurrency/quotes/latest", nil)
	if err != nil {
		return fmt.Sprintf(errStr, err)
	}

	q := url.Values{}
	q.Add("symbol", strings.Join(symbols, ","))

	req.Header.Set("Accepts", "application/json")
	req.Header.Add("X-CMC_PRO_API_KEY", "")
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Sprintf(errStr, err)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Sprintf(errStr, err)
	}

	var result map[string]any
	err = json.Unmarshal(respBody, &result)
	if err != nil {
		return fmt.Sprintf(errStr, err)
	}

	data := result["data"].(map[string]any)

	var sb strings.Builder

	for _, v := range symbols {
		first := data[v].([]any)[0].(map[string]any)
		quote := first["quote"].(map[string]any)
		usd := quote["USD"].(map[string]any)
		price := usd["price"].(float64)
		percentChange24h := usd["percent_change_24h"].(float64)

		if v == "PEPE" {
			sb.WriteString(fmt.Sprintf("%s: $%.8f (%.2f%%) \n", v, price, percentChange24h))
		} else {
			sb.WriteString(fmt.Sprintf("%s: $%.2f (%.2f%%) \n", v, price, percentChange24h))
		}

	}

	return sb.String()
}
