package server

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

type UsdBrlJson struct {
	USDBRL struct {
		Bid string `json:"bid"`
	} `json:"USDBRL"`
}

func server() {
	http.HandleFunc("/cotacao", cotacaoHandler)
	http.ListenAndServe(":8080", nil)
}

func cotacaoHandler(w http.ResponseWriter, r *http.Request) {
	cotacaoRequest(w, r)
}

func cotacaoRequest(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), time.Millisecond*200)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var apiJson UsdBrlJson
	err = json.Unmarshal(body, &apiJson)
	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(apiJson)
}
