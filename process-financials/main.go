package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	// Open the file
	r := createFileReader()

	// Iterate through the records
	count := sortFileContents(r)
	fmt.Println(count)
}

func sortFileContents(r *csv.Reader) int {
	count := 0
	var (
		amazonTotal float64 = 0
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

		fmt.Printf("%s\n", transaction)
		amazonTotal = filterTransactions("AMZN", transaction, amazonTotal)
		amazonTotal = filterTransactions("Amazon", transaction, amazonTotal)
	}
	fmt.Printf("amazon: %.2f", amazonTotal)
	return count
}

func filterTransactions(filterString string, transaction []string, categorySpend float64) float64 {

	if strings.Contains(transaction[2], filterString) {
		transactionAmount, err := strconv.ParseFloat(transaction[5], 64)
		if err != nil {
			log.Fatal(err)
		}
		categorySpend = categorySpend + transactionAmount
	}

	return categorySpend
}

func createFileReader() *csv.Reader {
	fileToOpen := "Chase8873_Activity20200501_20210527_20210527.CSV"
	csvfile, err := os.Open(fileToOpen)
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}

	// Parse the file
	r := csv.NewReader(csvfile)
	//r := csv.NewReader(bufio.NewReader(csvfile))
	return r
}