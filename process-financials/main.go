package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	fileToOpen := "Chase8873_Activity20200501_20210527_20210527.CSV"
	reader := createFileReader(fileToOpen)

	// add categories to filter with here
	categories := []category{
		newCategory("amazonTotal", []string{"AMZN", "Amazon"}),
	}

	sortFileContents(reader, categories)
}

func sortFileContents(reader *csv.Reader, allCategories []category) {
	numberOfCategories := len(allCategories)

	for categoryIndex := 0; categoryIndex < numberOfCategories; categoryIndex++ {
		currentCategory := allCategories[categoryIndex]
		for {
	    	// Read each transaction from csv
	    	transaction, err := reader.Read()
	    	if err == io.EOF {
	    		break
	    	}
	    	if err != nil {
	    		log.Fatal(err)
	    	}

			// loop through each filter in the category
			numberOfFilters := len(currentCategory.filters)
	    	for filterIndex := 0; filterIndex < numberOfFilters; filterIndex++ {
	    		currentFilter := currentCategory.filters[filterIndex]
	    		compareTransactionToFilter(currentFilter, currentCategory.values, transaction)
			}
	    }

	    // left at end so it can add everything together
		printCategoryResults(currentCategory)
	}

}

func compareTransactionToFilter(filterString string, totalCategorySpend map[string]float64, transaction []string) map[string]float64 {

	transactionDescription := transaction[2]
	if strings.Contains(transactionDescription, filterString) {

		// Get transaction amount
		transactionAmount, err := strconv.ParseFloat(transaction[5], 64)
		if err != nil {
			log.Fatal(err)
		}

		const transactionLayout = "01/02/2006"
		const sortingLayout= "2006 01"
		// Get the date of the transaction
	    transactionDate, err := time.Parse(transactionLayout, transaction[0])
	    if err != nil {
	    	log.Fatal(err)
		}

		monthSpend := transactionDate.Format(sortingLayout)
		// Add transaction value to the monthly total
		totalCategorySpend[monthSpend] = totalCategorySpend[monthSpend] + transactionAmount
	}

	return totalCategorySpend
}

func createFileReader(file string) *csv.Reader {
	// Open the file
	csvfile, err := os.Open(file)
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}

	// Parse the file
	r := csv.NewReader(csvfile)
	return r
}