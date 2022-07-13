package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	dbase "vanilla-loan-application/db"
)

type Application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

// struct untuk flag
type config struct {
	port string
}

var cfg config

func main() {
	fmt.Println("Starting server...")

	// initialize database
	dbConnString := "postgres://root:123@postgres:5432/loan-app-db?sslmode=disable"
	db, err := dbase.DBConnect(dbConnString)
	if err != nil {
		fmt.Println(err.Error())
		panic("failed to create connection database")
	}

	// setup flag cli property, default value, desc
	flag.StringVar(&cfg.port, "port", ":4000", "HTTP network address")
	flag.Parse()

	// leveled log
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &Application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	//
	flag.VisitAll(func(f *flag.Flag) { // syntax flag.VisitAll untuk mengakses semua nilai flag
		infoLog.Printf("Starting server with flag %s, port used is %s", f.Name, f.Value)
	})

	routerMiddleware := dbase.DBMiddleware(app.routes(), db)

	srv := &http.Server{
		Addr:     cfg.port,
		ErrorLog: errorLog,
		Handler:  routerMiddleware,
	}

	srv.ListenAndServe()
}
