package main

import (
	"encoding/csv"
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
	// Open the file
	r := createFileReader()

	sortFileContents(r)
}

func sortFileContents(r *csv.Reader) {
	// Categories to sort transactions by
	var (
		amazonTotals = make(map[string]float64)
	)

	for {
		// Read each transaction from csv
		transaction, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		amazonTotals = filterTransactions("AMZN", amazonTotals, transaction)
		amazonTotals = filterTransactions("Amazon", amazonTotals, transaction)
	}

	printCategoryResults(amazonTotals)
}

func printCategoryResults(amazonTotals map[string]float64) {
	keys := make([]string, 0, len(amazonTotals))
	for k := range amazonTotals {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	fullTotal := 0.00
	for _, k := range keys {
		fmt.Printf("%s: %.2f\n", k, amazonTotals[k])
		fullTotal = fullTotal + amazonTotals[k]
	}
	fmt.Printf("full total: %.2f\n", fullTotal)
}

func filterTransactions(filterString string, totalCategorySpend map[string]float64, transaction []string) map[string]float64 {

	if strings.Contains(transaction[2], filterString) {
		transactionAmount, err := strconv.ParseFloat(transaction[5], 64)
		if err != nil {
			log.Fatal(err)
		}

		const transactionLayout = "01/02/2006"
		const sortingLayout= "2006 01"
	    transactionDate, err := time.Parse(transactionLayout, transaction[0])
	    if err != nil {
	    	log.Fatal(err)
		}

		monthSpend := transactionDate.Format(sortingLayout)
		totalCategorySpend[monthSpend] = totalCategorySpend[monthSpend] + transactionAmount
	    //if _, ok := totalCategorySpend[monthSpend]; ok {
		//	totalCategorySpend[monthSpend] = totalCategorySpend[monthSpend] + transactionAmount
		//} else {
		//	totalCategorySpend[monthSpend] = transactionAmount
		//}
	}

	return totalCategorySpend
}

func createFileReader() *csv.Reader {
	fileToOpen := "Chase8873_Activity20200501_20210527_20210527.CSV"
	csvfile, err := os.Open(fileToOpen)
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}

	// Parse the file
	r := csv.NewReader(csvfile)
	return r
}