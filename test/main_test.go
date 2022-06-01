package test

import (
	"database/sql"
	db "github.com/freedommmoto/test_simplebank/db/sqlc"
	"log"
	"os"
	"testing"

	tool "github.com/freedommmoto/test_simplebank/tool"
	_ "github.com/lib/pq"
)

var testQueries *db.Queries

func TestMain(m *testing.M) {

	config, err := tool.LoadConfig("../")
	//fmt.Printf("%v", config)

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testQueries = db.New(conn)

	os.Exit(m.Run())
}
