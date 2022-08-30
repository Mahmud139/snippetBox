package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)

	fileserver := http.FileServer(http.Dir("M:/Projects/snippetbox/ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileserver))

	log.Println("Starting server on localhost:8080")
	err := http.ListenAndServe("localhost:8080", mux)
	log.Fatal(err)
}

