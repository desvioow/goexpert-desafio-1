package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type UsdBrlJson struct {
	Bid string `json:"bid"`
}

func main() {

	apiJson, err := brlUsdApiRequest()
	if err != nil {
		log.Fatal(err)
	}

	err = saveBidToFile(apiJson.Bid)
	if err != nil {
		log.Fatal(err)
	}
}

func brlUsdApiRequest() (*UsdBrlJson, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*300)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var apiJson UsdBrlJson
	err = json.Unmarshal(body, &apiJson)
	if err != nil {
		return nil, err
	}

	return &apiJson, nil
}

func saveBidToFile(bid string) error {

	file, err := os.OpenFile("bids.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = fmt.Fprintln(file, "Dólar: {"+bid+"}")
	if err != nil {
		return err
	}

	return nil
}
