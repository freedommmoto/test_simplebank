package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	tool "github.com/freedommmoto/test_simplebank/tool"
	_ "github.com/lib/pq"
)

var testQueries *Queries

func TestMain(m *testing.M) {

	config, err := tool.LoadConfig("../../")
	fmt.Printf("%v", config)

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testQueries = New(conn)

	os.Exit(m.Run())
}
