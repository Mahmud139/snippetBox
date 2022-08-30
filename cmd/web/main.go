package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	addr := flag.String("addr", "localhost:8080", "HTTP Netword Address")
	flag.Parse()
	
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)

	fileserver := http.FileServer(http.Dir("M:/Projects/snippetbox/ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileserver))

	log.Printf("Starting server on %s\n",*addr)
	err := http.ListenAndServe(*addr, mux)
	log.Fatal(err)
}

