package main

import (
	"database/sql"
	"log"
	"net"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/manucher051299/crud/cmd/app"
	"github.com/manucher051299/crud/pkg/customers"
)

func main() {
	host := "0.0.0.0"
	port := "9999"
	dsn := "postgres://app:pass@localhost:5432/db"
	err := execute(host, port, dsn)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
}

func execute(host string, port string, dsn string) (err error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Print(err)
		os.Exit(1)
		return
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Print(err)
			return
		}
	}()
	mux := http.NewServeMux()
	customersSvs := customers.NewService(db)
	server := app.NewServer(mux, customersSvs)
	server.Init()
	srv := &http.Server{
		Addr:    net.JoinHostPort(host, port),
		Handler: server,
	}
	return srv.ListenAndServe()
}

/*

	ctx := context.Background()
	_, err = db.ExecContext(ctx, `
	CREATE TABLE IF NOT EXISTS customers(
		id BIGSERIAL PRIMARY KEY,
		name 	TEXT 		NOT NULL,
		phone 	TEXT 		NOT NULL UNIQUE,
		active 	BOOLEAN 	NOT NULL DEFAULT TRUE,
		created TIMESTAMP 	NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		log.Print(err)
		return
	}
	name, phone := "Vasya", "+992000000001"
	result, err := db.ExecContext(ctx, `
	INSERT INTO customers(name,phone) VALUES ($1,$2) ON CONFLICT(phone) DO UPDATE SET name=excluded.name;
	`, name, phone)
	if err != nil {
		log.Print(err)
		return
	}
	id := 1
	newName := "Vasiliy"
	result, err = db.ExecContext(ctx, `
	UPDATE customers SET name=$2 WHERE id=$1;
	`, id, newName)
	if err != nil {
		log.Print(err)
		return
	}
	log.Print(result.RowsAffected())
	log.Print(result.LastInsertId())
	customer := &Customer{}
	err = db.QueryRowContext(ctx, `
	SELECT id,name,phone,active,created FROM customers WHERE id=1;
	`).Scan(&customer.ID, &customer.Name, &customer.Phone, &customer.Active, &customer.Created)
	if err != nil {
		log.Print(err)
		return
	}
	log.Printf("%#v", customer)
*/
