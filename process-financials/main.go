package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

func main() {
	var fileToOpen string
	flag.StringVar(&fileToOpen, "f", "", "File to be processed")
	flag.Parse()

	// Make sure -f flag set
	if fileToOpen == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	reader := createFileReader(fileToOpen)

	// add categories to filter with here
	categories := []category{
		newCategory("all", []string{"Sale"}),
		//newCategory("amazonTotal", []string{"AMZN", "Amazon"}),
	}

	sortFileContents(reader, categories)
}

func sortFileContents(reader *csv.Reader, allCategories []category) {
	numberOfCategories := len(allCategories)

	for {
		// Read each transaction from csv
		transaction, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		for categoryIndex := 0; categoryIndex < numberOfCategories; categoryIndex++ {
			currentCategory := allCategories[categoryIndex]

		    // loop through each filter in the category
	    	numberOfFilters := len(currentCategory.filters)
	    	for filterIndex := 0; filterIndex < numberOfFilters; filterIndex++ {

	    		transactionAdded := false
	    		currentFilter := currentCategory.filters[filterIndex]

	    		compareTransactionToFilter(currentFilter, currentCategory.values, transaction, &transactionAdded)

	    		if transactionAdded {
	    			break
	    		}
	    	}
	    }
	}

	// left at end so it can add everything together
	printCategoryResults(allCategories, numberOfCategories)

}

func compareTransactionToFilter(filterString string, totalCategorySpend map[string]float64, transaction []string, transactionAdded *bool) map[string]float64 {

	transactionDescription := transaction[2]
	transactionCategory := transaction[3]
	transactionType := transaction[4]

	if strings.Contains(transactionDescription, filterString) ||
		strings.Contains(transactionCategory, filterString) ||
		strings.Contains(transactionType, filterString) {

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
		*transactionAdded = true

		// to see exact transactions
		fmt.Printf("%s, %.2f, %s,\n", transactionDate.Format(transactionLayout), transactionAmount, transactionDescription)
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

func printCategoryResults(categories []category, numberOfCategories int) {
	fullTotal := 0.00
	for categoryIndex := 0; categoryIndex < numberOfCategories; categoryIndex++ {
		currentCategory := categories[categoryIndex]
		values := currentCategory.values
		name := currentCategory.name

		// make new slice of strings so keys can be sorted
		keys := make([]string, 0, len(values))
		for key := range values {
			keys = append(keys, key)
		}
		sort.Strings(keys)

		categoryTotal := 0.00
		// print sorted values from above step
		for _, k := range keys {
			fmt.Printf("%s %s: %.2f\n", name, k, values[k])
			categoryTotal = categoryTotal + values[k]
		}
		fmt.Printf("%s total: %.2f\n", name, categoryTotal)
		fullTotal = fullTotal + categoryTotal
	}
	fmt.Printf("full total: %.2f\n", fullTotal)
}

// use the function, not the struct directly
type category struct {
	name string
	filters []string
	values map[string]float64
}

func newCategory(name string, filters []string) category {
	return category{
		name: name,
		filters: filters,
		values: make(map[string]float64),
	}
}
