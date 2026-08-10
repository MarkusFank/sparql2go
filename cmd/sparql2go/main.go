package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/MarkusFank/sparql2go/internal/app"
)

var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

func main() {

	fmt.Printf("Welcome to sparql2go! You are using version %s, commit %s, built %s\n", version, commit, date)

	inputFile := flag.String("input", "", "RDF file to read")
	port := flag.Int("Port", 4711, "Port of the local web server that is being started")

	flag.Parse()

	err := app.Run(*inputFile, *port)

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
