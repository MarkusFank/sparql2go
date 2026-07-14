package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/MarkusFank/sparql2go/internal/app"
)

func main() {

	fmt.Println("Welcome to sparql2go!")
	inputFile := flag.String("input", "", "RDF file to read")
	port := flag.Int("Port", 4711, "Port of the local web server that is being started")

	flag.Parse()

	err := app.Run(*inputFile, *port)

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
