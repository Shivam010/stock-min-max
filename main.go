package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
)

// Response Object
type response struct {
	Error string                 `json:"error"`
	Data  interface{} `json:"data"`
}

func (r *response) Json() []byte {
	if r == nil {
		return []byte{}
	}
	b, _ := json.Marshal(r)
	return b
}

func (r *response) String() string {
	return string(r.Json())
}

// Return returns and adds the corresponding Valuer to the ResponseWriter
func _return(w http.ResponseWriter, data *response) {
	w.Header().Set("Content-Type", "application/json")
	Ignore(fmt.Fprintln(w, data))
}

type handler struct{}

func (*handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, "/api/v1") {
		_return(w, &response{Error: "404 page not found"})
	}
	commodity := r.URL.Query().Get("commodity")
	if commodity == "" {
		_return(w, &response{Error: "commodity must be specified"})
		return
	}
	_return(w, process(commodity))
}

func main() {
	// parse all flags set in `init`
	flag.Parse()

	fmt.Printf("\nStarting StockMinMax Server at port(:%v)...\n\n", PORT)

	if err := http.ListenAndServe(":"+PORT, &handler{}); err != nil {
		log.Println("\n\nPanics", err)
	}
}
