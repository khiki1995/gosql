package main

import (
	"net/http"
	"net"
	"os"
	"log"
	"database/sql"
	"github.com/khiki1995/gosql/cmd/app"
	"github.com/khiki1995/gosql/pkg/customers"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func main() {
	host := "0.0.0.0"
	port := "9999"
	dsn := "postgres://app:pass@localhost:5432/db"

	if err := execute(host, port, dsn); err != nil {
		log.Print(err)
		os.Exit(1)
	}
}

func execute(host string, port string, dsn string) (err error) {
	db, err := sql.Open("pgx",dsn)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := db.Close(); cerr != nil {
			err = cerr
			return
		}
		log.Println(err)
	}()
	mux := http.NewServeMux()
	customersSvc := customers.NewService(db)
	server := app.NewServer(mux, customersSvc)
	server.Init()

	srv := &http.Server{
		Addr: net.JoinHostPort(host, port),
		Handler: server,
	}
	return srv.ListenAndServe()
}