package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type application struct {
	errorLog *log.Logger
	infoLog *log.Logger
}

func main() {
	addr := flag.String("addr", "localhost:8080", "HTTP Netword Address")
	flag.Parse()
	
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application{
		errorLog: errorLog,
		infoLog: infoLog,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)

	fileserver := http.FileServer(http.Dir("M:/Projects/snippetbox/ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileserver))

	server := &http.Server{
		Addr: *addr,
		ErrorLog: errorLog,
		Handler: mux,
	}

	infoLog.Printf("Starting server on %s\n",*addr)
	err := server.ListenAndServe()
	errorLog.Fatal(err)
}

