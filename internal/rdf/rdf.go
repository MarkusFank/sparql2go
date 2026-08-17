package rdf

import (
	"fmt"
	"os"
	"path/filepath"

	rdflibgo "github.com/tggo/goRDFlib"
	"github.com/tggo/goRDFlib/nt"
	"github.com/tggo/goRDFlib/sparql"
	"github.com/tggo/goRDFlib/turtle"
)

var graph *rdflibgo.Graph
var rdfFile string

var lastQueryResult *sparql.Result

func Graph() *rdflibgo.Graph {
	return graph
}

func File() string {
	return rdfFile
}

func LastQueryResult() *sparql.Result {
	return lastQueryResult
}

func Init(inputFile string) error {
	file, err := os.Open(inputFile)

	if err != nil {
		return err
	}

	ext := filepath.Ext(inputFile)
	graph = rdflibgo.NewGraph()
	switch ext {
	case ".nt":
		if err = nt.Parse(graph, file, nt.WithUnboundedLines()); err != nil {
			return err
		}
	case ".ttl":
		if err = turtle.Parse(graph, file); err != nil {
			return err
		}
	default:
		return fmt.Errorf("Unable to parse file with extension %q", ext)
	}

	rdfFile = inputFile
	return nil
}

func Query(queryText string) (*sparql.Result, error) {
	res, err := sparql.Query(graph, queryText)

	if err != nil {
		lastQueryResult = nil
	} else if res != nil {
		lastQueryResult = res
	}

	return res, err
}
