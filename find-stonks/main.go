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

	incomeStatement := getInfoFromApi("INCOME_STATEMENT", stockSymbol, apiKey)
	incomeStatementObject := convertIncomeStatementToObject(incomeStatement)
	projectedRevenue := calculateAverageYearlyRevenueGrowth(incomeStatementObject)
	fmt.Printf("proj revenue: %f\n", projectedRevenue)
	fmt.Println("")

	balanceSheet := getInfoFromApi("BALANCE_SHEET", stockSymbol, apiKey)
	balanceSheetObject := convertBalanceSheetToObject(balanceSheet)
	projectedShares := calculateAverageYearlyShareGrowth(balanceSheetObject)
	fmt.Printf("proj share count: %f\n", projectedShares)
	fmt.Println("")

	averageMargin := calculateAverageMargin(incomeStatementObject)
	fmt.Printf("avg margin: %f\n", averageMargin)
	fmt.Println("")

	companyOverview := getInfoFromApi("OVERVIEW", stockSymbol, apiKey)
	companyOverviewObject := convertCompanyOverviewToObject(companyOverview)
	currentPEString := companyOverviewObject["PERatio"]
	currentPE, err := strconv.ParseFloat(currentPEString, 64)
	if err != nil {
		log.Fatal(err)
	}
	adjustedPE := currentPE * .8
	fmt.Printf("current PE: %f\n", currentPE)
	fmt.Printf("adjusted PE: %f\n", adjustedPE)
	fmt.Println("")

	priceEstimate := adjustedPE * averageMargin * projectedRevenue / projectedShares
	fmt.Printf("priceEstimate: %.2f\n", priceEstimate)


	fmt.Printf("15 percent price: %.2f\n", priceEstimate / math.Pow(1.15, 5))
	fmt.Printf("20 percent price: %.2f\n", priceEstimate / math.Pow(1.20, 5))
	fmt.Printf("25 percent price: %.2f\n", priceEstimate / math.Pow(1.25, 5))
	fmt.Printf("30 percent price: %.2f\n", priceEstimate / math.Pow(1.30, 5))

}

func calculateAverageMargin(incomeStatementObject incomeStatementResponse) float64 {
	incomeAnnualStatements := incomeStatementObject.AnnualReports
	var totalMargin float64
	yearsOfData := len(incomeAnnualStatements)

	// its conceivable there could be a bug if the AnnualReports are returned out order
	for i := 0; i < yearsOfData - 1; i++ {
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
		totalMargin = totalMargin + margin
	}
	// yearsOfData - 1 because that's the number of comparisons that can be made
	averageMargin := totalMargin / float64(yearsOfData - 1)

	return averageMargin
}

func calculateAverageYearlyShareGrowth(balanceSheetObject balanceSheetResponse) float64 {

	balanceSheetAnnuals := balanceSheetObject.AnnualReports
	var totalShareGrowthPercent float64
	yearsOfData := len(balanceSheetAnnuals)

	// its conceivable there could be a bug if the AnnualReports are returned out order
	for i := 0; i < yearsOfData - 1; i++ {
		shareCountString := balanceSheetAnnuals[i]["commonStockSharesOutstanding"]
		previousShareCountString := balanceSheetAnnuals[i+1]["commonStockSharesOutstanding"]

		currentShareCount, err := strconv.Atoi(shareCountString)
		if err != nil {
			log.Fatal(err)
		}
		previousShareCount, err := strconv.Atoi(previousShareCountString)
		if err != nil {
			log.Fatal(err)
		}

		// For some reason this did not want to work with int, I kept getting 0 and had to make them all floats
		yearlyShareGrowthPercent := float64(currentShareCount- previousShareCount) / float64(previousShareCount)
		totalShareGrowthPercent = totalShareGrowthPercent + yearlyShareGrowthPercent
	}
	// yearsOfData - 1 because that's the number of comparisons that can be made
	averageShareGrowthPercent := totalShareGrowthPercent / float64(yearsOfData - 1)

	return fiveYearShareProjection(averageShareGrowthPercent, balanceSheetAnnuals)
}

func fiveYearShareProjection(averageShareGrowthPercent float64, balanceSheetAnnuals []map[string]string) float64 {
	// multiply by 1.1 to be conservative
	adjustedAverageShareGrowthPercent := averageShareGrowthPercent * 1.1
	projectionMultiple := math.Pow(1+adjustedAverageShareGrowthPercent, 5)
	latestShareCount, err := strconv.ParseFloat(balanceSheetAnnuals[0]["commonStockSharesOutstanding"], 64)
	if err != nil {
		log.Fatal(err)
	}
	sharesInFiveYears := projectionMultiple * latestShareCount / 1000000

	fmt.Printf("current share count: %f\n", latestShareCount / 1000000)
	return sharesInFiveYears
}

func calculateAverageYearlyRevenueGrowth(incomeStatementObject incomeStatementResponse) float64 {

	incomeAnnualStatements := incomeStatementObject.AnnualReports
	var totalRevenueGrowthPercent float64
	yearsOfData := len(incomeAnnualStatements)

	// its conceivable there could be a bug if the AnnualReports are returned out order
	for i := 0; i < yearsOfData - 1; i++ {
		yearlyRevenueString := incomeAnnualStatements[i]["totalRevenue"]
		previousYearlyRevenueString := incomeAnnualStatements[i+1]["totalRevenue"]

		yearlyRevenue, err := strconv.Atoi(yearlyRevenueString)
		if err != nil {
			log.Fatal(err)
		}
		previousYearlyRevenue, err := strconv.Atoi(previousYearlyRevenueString)
		if err != nil {
			log.Fatal(err)
		}

		// For some reason this did not want to work with int, I kept getting 0 and had to make them all floats
		yearlyRevenueGrowthPercent := float64(yearlyRevenue - previousYearlyRevenue) / float64(previousYearlyRevenue)
		totalRevenueGrowthPercent = totalRevenueGrowthPercent + yearlyRevenueGrowthPercent
	}
	// yearsOfData - 1 because that's the number of comparisons that can be made
	averageRevenueGrowthPercent := totalRevenueGrowthPercent / float64(yearsOfData - 1)

	return fiveYearRevenueProjection(averageRevenueGrowthPercent, incomeAnnualStatements)
}

func fiveYearRevenueProjection(averageRevenueGrowthPercent float64, incomeAnnualStatements []map[string]string) float64 {
	// multiply by .9 to be conservative
	adjustedAverageRevenueGrowthPercent := averageRevenueGrowthPercent * .9
	projectionMultiple := math.Pow(1+adjustedAverageRevenueGrowthPercent, 5)
	latestRevenueCount, err := strconv.ParseFloat(incomeAnnualStatements[0]["totalRevenue"], 64)
	if err != nil {
		log.Fatal(err)
	}
	revenueInFiveYears := projectionMultiple * latestRevenueCount / 1000000

	fmt.Printf("latest: %f\n", latestRevenueCount / 1000000)
	fmt.Printf("adj revenue growth: %f\n", adjustedAverageRevenueGrowthPercent)
	return revenueInFiveYears
}

func convertIncomeStatementToObject(apiOutput []byte) incomeStatementResponse {
	var newObject incomeStatementResponse

	if err := json.Unmarshal(apiOutput, &newObject); err != nil {
		log.Fatal(err)
	}

	return newObject
}

func convertBalanceSheetToObject(apiOutput []byte) balanceSheetResponse {
	var newObject balanceSheetResponse

	if err := json.Unmarshal(apiOutput, &newObject); err != nil {
		log.Fatal(err)
	}

	return newObject
}

func convertCompanyOverviewToObject(apiOutput []byte) map[string]string {
	var newObject map[string]string

	if err := json.Unmarshal(apiOutput, &newObject); err != nil {
		log.Fatal(err)
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
