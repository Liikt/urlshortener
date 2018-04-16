package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/liikt/urlshort/handler"
)

var (
	ymlpath  string
	boltpath string
)

func main() {
	flag.StringVar(&ymlpath, "ymlpath", "storage/map.yml", "Path to a YAML file containing shortened URLs")
	flag.StringVar(&boltpath, "boltpath", "storage/bolt.db", "Path to a BoltDB File containing shortened URLs")
	flag.Parse()

	mux := defaultMux()

	if yaml, err := ioutil.ReadFile(ymlpath); err == nil {
		err = handler.YAMLHandler([]byte(yaml), mux)
		if err != nil {
			panic(err)
		}
	}

	err := handler.BoltHandler(boltpath, mux)
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", mux)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	mux.HandleFunc("/api/show", show)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}

func show(w http.ResponseWriter, _ *http.Request) {
	content, _ := handler.DBContent(boltpath)
	for key, val := range content {
		w.Write([]byte(key + ":" + val + "\n"))
	}
}
