package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

func main() {
	addr := flag.String("addr", "localhost:8080", "HTTP Netword Address")
	flag.Parse()
	
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)

	fileserver := http.FileServer(http.Dir("M:/Projects/snippetbox/ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileserver))

	infoLog.Printf("Starting server on %s\n",*addr)
	err := http.ListenAndServe(*addr, mux)
	errorLog.Fatal(err)
}

