package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	stockSymbol := processArgs()

	apiKey := os.Getenv("ALPHA_ADVANTAGE_PW")
	stockFunction := "BALANCE_SHEET"
	url := fmt.Sprintf("https://www.alphavantage.co/query?function=%s&symbol=%s&apikey=%s", stockFunction, stockSymbol, apiKey)

	req, _ := http.NewRequest("GET", url, nil)
	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(res)
	fmt.Println(string(body))

}