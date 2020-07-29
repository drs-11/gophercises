package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"./handler"
)

func main() {
	YAMLFile := flag.String("yaml", "mappings.yml", "YAML file where the url mappings are located")
	mux := defaultMux()

	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}

	mapHandler := handler.MapHandler(pathsToUrls, mux)

	file, err := os.Open(*YAMLFile)
	if err != nil {
		log.Fatal(err)
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()

	fmt.Println("Running server on localhost, port 8080")
	redirectHandler, err := handler.YAMLHandler(data, mapHandler)
	if err != nil {
		log.Fatal(err)
	}

	http.ListenAndServe(":8080", redirectHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", root)
	return mux
}

func root(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "You have accessed the root directory of the server.")
}
