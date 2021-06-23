package main

import (
	"flag"
	"os"
)

func processArgs() string {
	var stockSymbol string

	flag.StringVar(&stockSymbol, "s", "", "Stock to process and see if it's a stonk")
	flag.Parse()

	// Make sure -s flag set
	if stockSymbol == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	return stockSymbol
}
