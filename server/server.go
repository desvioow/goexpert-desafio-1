package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type UsdBrlJson struct {
	USDBRL struct {
		Bid string `json:"bid"`
	} `json:"USDBRL"`
}

func main() {
	http.HandleFunc("/cotacao", cotacaoHandler)
	http.ListenAndServe(":8080", nil)
}

func cotacaoHandler(w http.ResponseWriter, r *http.Request) {

	apiJson, err := cotacaoRequest(w, r)
	if err != nil {
		log.Fatal(err)
	}

	db, err := retrieveDb()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	bid, err := strconv.ParseFloat(apiJson.USDBRL.Bid, 64)
	if err != nil {
		log.Fatal(err)
	}
	err = insertBidIntoDb(db, bid)
	if err != nil {
		log.Fatal(err)
	}

}

func cotacaoRequest(w http.ResponseWriter, r *http.Request) (*UsdBrlJson, error) {

	ctx, cancel := context.WithTimeout(r.Context(), time.Millisecond*200)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
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

	json.NewEncoder(w).Encode(apiJson)
	return &apiJson, nil
}

func retrieveDb() (*sql.DB, error) {

	db, err := sql.Open("sqlite3", "sqlite.db")
	if err != nil {
		return nil, err
	}
	create_query := `CREATE TABLE IF NOT EXISTS usdxbrl (id INTEGER PRIMARY KEY AUTOINCREMENT, bid REAL);`

	_, err = db.Exec(create_query)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func insertBidIntoDb(db *sql.DB, bid float64) error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
	defer cancel()
	insert_query := `INSERT INTO usdxbrl (bid) values(?) `

	stmt, err := db.Prepare(insert_query)
	defer stmt.Close()
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, bid)
	if err != nil {
		return err
	}

	return nil
}
