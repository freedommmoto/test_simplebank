package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"

	"github.com/freedommmoto/test_simplebank/api"
	db "github.com/freedommmoto/test_simplebank/db/sqlc"
	tool "github.com/freedommmoto/test_simplebank/tool"
)

var mainQueries *db.Queries

func connectDB() {
	config, err := tool.LoadConfig(".")
	fmt.Printf("%v", config)
	if err != nil {
		log.Fatal("cannot load config file:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
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

// func main() {

// 	http.HandleFunc("/hello", hello)
// 	http.HandleFunc("/headers", headers)
// 	connectDB() // if don't call this line before Listen is will get error invalid memory address or nil pointer dereference
// 	http.ListenAndServe(":8082", nil)
// }

func main() {

	config, err := tool.LoadConfig(".")
	fmt.Printf("%v", config)
	if err != nil {
		log.Fatal("cannot load config file:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	store := db.NewStore(conn)
	server, nil := api.NewServer(config, store)
	if err != nil {
		log.Fatal("can't NewServer", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("can't start server with gin", err)
	}

}
