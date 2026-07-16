package app

import (
	"errors"
	"os"
	"strings"

	"github.com/MarkusFank/sparql2go/internal/rdf"
	"github.com/MarkusFank/sparql2go/internal/webserver"
)

func Run(inputFile string, port int) error {

	if len(strings.TrimSpace(inputFile)) == 0 {
		return errors.New("Input file name not specified")
	}

	_, err := os.Stat(inputFile)

	if err != nil {
		return err
	}

	err = rdf.Init(inputFile)

	if err != nil {
		return err
	}

	err = webserver.Run(port, inputFile)

	if err != nil {
		return err
	}

	return nil
}
