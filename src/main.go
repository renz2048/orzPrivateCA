package main

import (
	"fmt"
	"log"
	"net/http"
)

type Engine struct{}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/":
		fmt.Fprintf(w, "URL.path = %q\n", req.URL.Path)
	case "/rootca":
		for k, v := range req.Header {
			fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
		}
	default:
		fmt.Fprintf(w, "404 NOT FOUND: %s\n", req.URL)
	}
}

func main() {
	/*
		http.HandleFunc("/", indexHandler)
		http.HandleFunc("/rootca", rootCAHandler)
		http.HandleFunc("/subca", subCAHandler)
		log.Fatal(http.ListenAndServeTLS(":443",
			"../certificate/www.example.com.crt",
			"../certificate/www.example.com.key.unsecure",
			nil))
	*/
	engine := new(Engine)
	log.Fatal(http.ListenAndServeTLS(":443",
		"../certificate/www.example.com.crt",
		"../certificate/www.example.com.key.unsecure",
		engine))
}

func indexHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
}

func rootCAHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Root CA\n")
}

func subCAHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Subordinate CA\n")
}
