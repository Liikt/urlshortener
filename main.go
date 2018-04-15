package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gophercises/urlshort"
)

func main() {
	ymlpath := flag.String("ymlpath", "../storage/map.yml", "Path to a YAML file containing shortened URLs")
	boltpath := flag.String("boltpath", "../storage/bolt.db", "Path to a BoltDB File containing shortened URLs")
	flag.Parse()

	mux := defaultMux()

	if yaml, err := ioutil.ReadFile(*ymlpath); err == nil {
		err = urlshort.YAMLHandler([]byte(yaml), mux)
		if err != nil {
			panic(err)
		}
	}

	err := urlshort.BoltHandler(*boltpath, mux)
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", mux)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
