package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"
)

type FinancialsResponse struct {
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
	//currentPrice := GetInfoFromApi("GLOBAL_QUOTE", stockSymbol, apiKey)
	//fmt.Printf("current price: %f\n", currentPrice)

	incomeStatement := GetInfoFromApi("INCOME_STATEMENT", stockSymbol, apiKey)
	incomeStatementObject := ConvertFinancialsResponseToObject(incomeStatement)
	projectedRevenue := CalculateAverageYearlyMetricGrowth("totalRevenue", incomeStatementObject)

	balanceSheet := GetInfoFromApi("BALANCE_SHEET", stockSymbol, apiKey)
	balanceSheetObject := ConvertFinancialsResponseToObject(balanceSheet)
	projectedShares := CalculateAverageYearlyMetricGrowth("commonStockSharesOutstanding", balanceSheetObject)

	averageMargin := CalculateAverageMargin(incomeStatementObject)

	companyOverview := GetInfoFromApi("OVERVIEW", stockSymbol, apiKey)
	companyOverviewObject := ConvertCompanyOverviewToObject(companyOverview)
	adjustedPE := AdjustCurrentPE(companyOverviewObject)

	outputPriceEstimates(adjustedPE, averageMargin, projectedRevenue, projectedShares)
}

func outputPriceEstimates(adjustedPE float64, averageMargin float64, projectedRevenue float64, projectedShares float64) {
	priceEstimate := adjustedPE * averageMargin * projectedRevenue / projectedShares
	fmt.Println("")
	fmt.Printf("5 yr priceEstimate: %.2f\n", priceEstimate)
	fmt.Printf("15 percent price: %.2f\n", priceEstimate/math.Pow(1.15, 5))
	fmt.Printf("20 percent price: %.2f\n", priceEstimate/math.Pow(1.20, 5))
	fmt.Printf("25 percent price: %.2f\n", priceEstimate/math.Pow(1.25, 5))
	fmt.Printf("30 percent price: %.2f\n", priceEstimate/math.Pow(1.30, 5))
}

func AdjustCurrentPE(companyOverviewObject map[string]string) float64 {
	currentPEString := companyOverviewObject["PERatio"]
	currentPE, err := strconv.ParseFloat(currentPEString, 64)
	if err != nil {
		log.Fatal(err)
	}

	adjustedPE := currentPE * .8

	fmt.Println("")
	fmt.Printf("current PE: %f\n", currentPE)
	fmt.Printf("adjusted PE: %f\n", adjustedPE)

	return adjustedPE
}

func CalculateAverageMargin(incomeStatementObject FinancialsResponse) float64 {
	incomeAnnualStatements := incomeStatementObject.AnnualReports
	var totalMargin float64
	yearsOfData := len(incomeAnnualStatements)

	// its conceivable there could be a bug if the AnnualReports are returned out order
	for i := 0; i < yearsOfData-1; i++ {
		yearlyNetIncomeString := incomeAnnualStatements[i]["netIncome"]
		yearlyRevenueString := incomeAnnualStatements[i]["totalRevenue"]

		yearlyNetIncome, err := strconv.Atoi(yearlyNetIncomeString)
		if err != nil {
			log.Fatal(err)
		}
		yearlyRevenue, err := strconv.Atoi(yearlyRevenueString)
		if err != nil {
			log.Fatal(err)
		}

		// For some reason this did not want to work with int, I kept getting 0 and had to make them all floats
		margin := float64(yearlyNetIncome) / float64(yearlyRevenue)
		totalMargin += margin
	}
	// yearsOfData - 1 because that's the number of comparisons that can be made
	averageMargin := totalMargin / float64(yearsOfData-1)

	fmt.Println("")
	fmt.Printf("avg margin: %f\n", averageMargin)

	return averageMargin
}

func CalculateAverageYearlyMetricGrowth(metric string, financialsObject FinancialsResponse) float64 {

	annualReports := financialsObject.AnnualReports
	var totalMetricGrowthPercent float64
	yearsOfData := len(annualReports)

	// its conceivable there could be a bug if the AnnualReports are returned out order
	for i := 0; i < yearsOfData-1; i++ {
		currentMetricString := annualReports[i][metric]
		previousMetricString := annualReports[i+1][metric]

		currentMetric, err := strconv.Atoi(currentMetricString)
		if err != nil {
			log.Fatal(err)
		}
		previousMetric, err := strconv.Atoi(previousMetricString)
		if err != nil {
			log.Fatal(err)
		}

		// For some reason this did not want to work with int, I kept getting 0 and had to make them all floats
		yearlyMetricGrowthPercent := float64(currentMetric-previousMetric) / float64(previousMetric)
		totalMetricGrowthPercent += yearlyMetricGrowthPercent
	}
	numberOfComparisions := float64(yearsOfData - 1)
	averageMetricGrowthPercent := totalMetricGrowthPercent / numberOfComparisions

	fmt.Println("")
	fmt.Printf("%s years of data: %v\n", metric, yearsOfData)

	return FiveYearMetricProjection(metric, averageMetricGrowthPercent, annualReports)
}

func FiveYearMetricProjection(metric string, percentChange float64, financialsObject []map[string]string) float64 {
	var adjustedPercentChange float64
	switch metric {
	case "totalRevenue":
		// multiply by .9 to assume less revenue
		adjustedPercentChange = percentChange * .9
	case "commonStockSharesOutstanding":
		// multiply by 1.1 to assume more shares outstanding
		adjustedPercentChange = percentChange * 1.1
	}

	projectionMultiple := math.Pow(1+adjustedPercentChange, 5)
	latestMetricValue, err := strconv.ParseFloat(financialsObject[0][metric], 64)
	if err != nil {
		log.Fatal(err)
	}

	// divide by 1 mil for easier numbers
	adjustedMetric := latestMetricValue / 1000000
	metricInFiveYears := projectionMultiple * adjustedMetric

	fmt.Printf("%s latest: %.3f\n", metric, adjustedMetric)
	fmt.Printf("yoy percent change: %.3f\n", adjustedPercentChange)
	fmt.Printf("5 year projection: %.3f\n", metricInFiveYears)
	return metricInFiveYears
}

func ConvertFinancialsResponseToObject(apiOutput []byte) FinancialsResponse {
	var newObject FinancialsResponse

	if err := json.Unmarshal(apiOutput, &newObject); err != nil {
		log.Fatal(err)
	}

	return newObject
}

func ConvertCompanyOverviewToObject(apiOutput []byte) map[string]string {
	var newObject map[string]string

	if err := json.Unmarshal(apiOutput, &newObject); err != nil {
		log.Fatal(err)
	}

	return newObject
}

func GetInfoFromApi(stockFunction string, stockSymbol string, apiKey string) []byte {
	url := fmt.Sprintf("https://www.alphavantage.co/query?function=%s&symbol=%s&apikey=%s", stockFunction, stockSymbol, apiKey)

	req, _ := http.NewRequest("GET", url, nil)
	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	return body
}
