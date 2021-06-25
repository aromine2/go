package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

type incomeStatementResponse struct {
	Symbol           string
	AnnualReports    []map[string]string
	QuarterlyReports []map[string]string
}

type balanceSheetResponse struct {
	Symbol           string
	AnnualReports    []map[string]string
	QuarterlyReports []map[string]string
}

type yearlyResults struct {
	Year                 time.Time
	RevenueGrowthPercent float64
	SharesGrowthPercent  float64
}

func main() {
	stockSymbol := processArgs()
	apiKey := os.Getenv("ALPHA_ADVANTAGE_PW")

	//balanceSheet := getInfoFromApi("BALANCE_SHEET", stockSymbol, apiKey)
	//balanceSheetObject := convertBalanceSheetToObject(balanceSheet)
	//
	//// Need to find a way to filter on the latest
	//fmt.Println(balanceSheetObject.AnnualReports[3]["longTermDebt"])

	averageRevenueGrowthPercent := calculateAverageYearlyRevenueGrowth(stockSymbol, apiKey)
	fmt.Println(averageRevenueGrowthPercent)

	// Get revenue growth rate over the last 4-5 years (%Inc = (F - I) / I)
	// Average that, then project 5 years out. Then take 90% of that to be conservative
	// Get shares growth rate over the last several years
	// Avg and project. Multiply by 1.1 to be conservative
	// Take average net profit to the projected revenue as a percent
	// Take the current P/E and multiply by .75 to be conservative
	// Return price estimate 5 years out
	// Return "today's stock price to buy at"; 15%, 20%, 25%, 30%
}

func calculateAverageYearlyRevenueGrowth(stockSymbol string, apiKey string) float64 {
	incomeStatement := getInfoFromApi("INCOME_STATEMENT", stockSymbol, apiKey)
	incomeStatementObject := convertIncomeStatementToObject(incomeStatement)

	incomeAnnualStatements := incomeStatementObject.AnnualReports
	var totalRevenueGrowthPercent float64
	yearsOfData := len(incomeAnnualStatements)

	// its conceivable there could be a bug if the AnnualReports are returned out order
	for i := 0; i < yearsOfData - 1; i++ {
		yearlyRevenueString := incomeAnnualStatements[i]["totalRevenue"]
		previousYearlyRevenueString := incomeAnnualStatements[i+1]["totalRevenue"]

		yearlyRevenue, err := strconv.Atoi(yearlyRevenueString)
		if err != nil {
			panic(err)
		}
		previousYearlyRevenue, err := strconv.Atoi(previousYearlyRevenueString)
		if err != nil {
			panic(err)
		}

		// For some reason this did not want to work with int, I kept getting 0 and had to make them all floats
		yearlyRevenueGrowthPercent := float64(yearlyRevenue - previousYearlyRevenue) / float64(previousYearlyRevenue)
		totalRevenueGrowthPercent = totalRevenueGrowthPercent + yearlyRevenueGrowthPercent
	}
	// yearsOfData - 1 because that's the number of comparisons that can be made
	averageRevenueGrowthPercent := totalRevenueGrowthPercent / float64(yearsOfData - 1)

	return averageRevenueGrowthPercent
}

func convertIncomeStatementToObject(apiOutput []byte) incomeStatementResponse {
	var newObject incomeStatementResponse

	if err := json.Unmarshal(apiOutput, &newObject); err != nil {
		panic(err)
	}

	return newObject
}

func convertBalanceSheetToObject(apiOutput []byte) balanceSheetResponse {
	var newObject balanceSheetResponse

	if err := json.Unmarshal(apiOutput, &newObject); err != nil {
		panic(err)
	}

	return newObject
}

func getInfoFromApi(stockFunction string, stockSymbol string, apiKey string) []byte {
	url := fmt.Sprintf("https://www.alphavantage.co/query?function=%s&symbol=%s&apikey=%s", stockFunction, stockSymbol, apiKey)

	req, _ := http.NewRequest("GET", url, nil)
	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	return body
}
