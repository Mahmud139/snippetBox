package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	errorLog *log.Logger
	infoLog *log.Logger
}

func main() {
	addr := flag.String("addr", "localhost:8080", "HTTP Netword Address")
	dsn := flag.String("dsn", "web:mahmud@/snippetbox?parseTime=true", "MySQL data source name")
	flag.Parse()
	
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application{
		errorLog: errorLog,
		infoLog: infoLog,
	}

	db, err := OpenDB(*dsn)																																														
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	server := &http.Server{
		Addr: *addr,
		ErrorLog: errorLog,
		Handler: app.routes(),
	}

	infoLog.Printf("Starting server on %s\n",*addr)
	err = server.ListenAndServe()
	errorLog.Fatal(err)
}

func OpenDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

