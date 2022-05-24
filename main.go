package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"

	db "github.com/mgo/banksystem/db/sqlc"
	"github.com/mgo/banksystem/tool"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@postgres-docker-service:5432/simple_bank?sslmode=disable"
)

var mainQueries *db.Queries

func connectDB() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	mainQueries = db.New(conn)
}

func hello(w http.ResponseWriter, req *http.Request) {

	fmt.Fprintf(w, "hello %v \n", tool.RandomOwner())

	arg := db.CreateCustomerParams{
		CustomerName: tool.RandomOwner(), // should used random data
		Balance:      tool.RandomMoney(),
		Currency:     tool.RandomCurrency(),
	}
	customerObject, err := mainQueries.CreateCustomer(context.Background(), arg)

	customerObjectAfterGet, err := mainQueries.GetCustomer(context.Background(), customerObject.ID)
	if err != nil {
		fmt.Fprintf(w, "errpr %v \n", err)
	}
	if int64(customerObjectAfterGet.ID) != int64(0) {
		fmt.Fprintf(w, "customer name %v \n", customerObject.CustomerName)
	}
}

func headers(w http.ResponseWriter, req *http.Request) {

	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func main() {

	http.HandleFunc("/hello", hello)
	http.HandleFunc("/headers", headers)
	connectDB() // if don't call this line before Listen is will get error invalid memory address or nil pointer dereference
	http.ListenAndServe(":8082", nil)
}
