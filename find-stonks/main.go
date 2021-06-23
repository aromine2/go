package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type balanceSheetResponse struct {
	Symbol    string
	AnnualReports   []map[string]string
	QuarterlyReports   []map[string]string
}

func main() {
	stockSymbol := processArgs()
	apiKey := os.Getenv("ALPHA_ADVANTAGE_PW")
	stockFunction := "BALANCE_SHEET"

	balanceSheet := getInfoFromApi(stockFunction, stockSymbol, apiKey)
	balanceSheetObject := processBalanceSheet(balanceSheet)

	// Need to find a way to filter on the latest
	fmt.Println(balanceSheetObject.AnnualReports[3]["longTermDebt"])
}

func processBalanceSheet(balanceSheet []byte) balanceSheetResponse {
	var stockBalanceSheet balanceSheetResponse

	err := json.Unmarshal(balanceSheet, &stockBalanceSheet)
	if err != nil {
		log.Println(err)
	}

	return stockBalanceSheet
}

func getInfoFromApi(stockFunction string, stockSymbol string, apiKey string) []byte {
	url := fmt.Sprintf("https://www.alphavantage.co/query?function=%s&symbol=%s&apikey=%s", stockFunction, stockSymbol, apiKey)

	req, _ := http.NewRequest("GET", url, nil)
	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	return body
}
