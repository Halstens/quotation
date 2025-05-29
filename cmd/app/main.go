package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/quotes/internal/config"
	"github.com/quotes/internal/database/postgress"
	"github.com/quotes/internal/repository"

	_ "github.com/lib/pq"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	quotes   *postgress.QuotesRepository
}

func main() {
	addr := flag.String("addr", ":4000", "Сетевой адрес HTTP")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime)

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Подключение к БД
	db, err := repository.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	db.SetMaxOpenConns(100)
	defer db.Close()

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		quotes:   &postgress.QuotesRepository{DB: db},
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Запуск сервера на %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}
