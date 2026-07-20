package webserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/MarkusFank/sparql2go/internal/rdf"
	"github.com/tggo/goRDFlib/sparql"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Vary", "Origin")

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func Run(port int, rdfFile string) error {

	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/query", handleQuery)

	mux.HandleFunc("GET /api/init", initEndpoint(rdfFile))

	fmt.Printf("Webserver is running on port %d. Open http://localhost:%d in your browser\n", port, port)
	handler := corsMiddleware(mux)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), handler)

	if err != nil {
		return err
	}

	return nil
}

type queryResult struct {
	Count  int                 `json:"count"`
	Result []map[string]string `json:"result"`
	Vars   []string            `json:"vars"`
}

func handleQuery(w http.ResponseWriter, r *http.Request) {

	// TODO use multipart form
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

	res := queryResult{Count: 0, Vars: queryRes.Vars, Result: []map[string]string{}}
	for _, row := range queryRes.Bindings {
		res.Count++
		mapObj := map[string]string{}
		for k, v := range row {
			mapObj[k] = v.String()
		}

		res.Result = append(res.Result, mapObj)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(res)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error while JSON encoding %v", err)
		return
	}
}

type initEndpointResponse struct {
	RdfFilePath string `json:"rdfFilePath"`
}

func initEndpoint(rdfFile string) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		res := initEndpointResponse{RdfFilePath: rdfFile}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
	}
}
