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
	Count  int                 `json:"count"`
	Result []map[string]string `json:"result"`
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

	res := queryResult{Count: 0, Result: []map[string]string{}}
	for i, row := range queryRes.Bindings {
		fmt.Printf("Row %d\n", i)
		res.Count++
		mapObj := map[string]string{}
		for k, v := range row {
			mapObj[k] = v.String()
		}

		res.Result = append(res.Result, mapObj)
	}

	fmt.Printf("This is the result: %v\n", res)

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(res)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error while JSON encoding %v", err)
		return
	}
}
