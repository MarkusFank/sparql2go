package rdf

import (
	"fmt"
	"os"
	"path/filepath"

	rdflibgo "github.com/tggo/goRDFlib"
	"github.com/tggo/goRDFlib/nt"
	"github.com/tggo/goRDFlib/turtle"
)

var Graph *rdflibgo.Graph

func Init(inputFile string) error {
	file, err := os.Open(inputFile)

	if err != nil {
		return err
	}

	ext := filepath.Ext(inputFile)
	Graph = rdflibgo.NewGraph()
	switch ext {
	case ".nt":
		if err = nt.Parse(Graph, file); err != nil {
			return err
		}
	case ".ttl":
		if err = turtle.Parse(Graph, file); err != nil {
			return err
		}
	default:
		return fmt.Errorf("Unable to parse file with extension %q", ext)
	}

	return nil
}
