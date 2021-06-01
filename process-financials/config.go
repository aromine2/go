package main

import (
    "flag"
    "os"
)

func processArgs() string {
    var fileToOpen string
    flag.StringVar(&fileToOpen, "f", "", "File to be processed")
    flag.Parse()

    // Make sure -f flag set
    if fileToOpen == "" {
        flag.PrintDefaults()
        os.Exit(1)
    }

    return fileToOpen
}
