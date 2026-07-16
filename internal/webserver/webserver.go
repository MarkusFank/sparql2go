package webserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/MarkusFank/sparql2go/internal/rdf"
	"github.com/tggo/goRDFlib/sparql"
)

func Run(port int, rdfFile string) error {

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", hello(rdfFile))

	mux.HandleFunc("POST /query", handleQuery)

	fmt.Printf("Webserver is running on port %d. Open http://localhost:%d in your browser\n", port, port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), mux)

	if err != nil {
		return err
	}

	return nil
}

func hello(rdfFile string) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, you are working with file %q", rdfFile)

		w.WriteHeader(200)
	}
}

type queryResult struct {
	count  int
	result []map[string]string
}

func handleQuery(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	query := r.Form.Get("query")

	if len(strings.TrimSpace(query)) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "No query provided")
		return
	}

	queryRes, err := sparql.Query(rdf.Graph, query)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unable to execute query %v", err)
		return
	}

	res := queryResult{count: 0, result: []map[string]string{}}
	for _, row := range queryRes.Bindings {
		res.count++
		mapObj := map[string]string{}
		for k, v := range row {
			mapObj[k] = v.String()
		}

		res.result = append(res.result, mapObj)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
