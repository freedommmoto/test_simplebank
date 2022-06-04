package test

import (
	"database/sql"
	"log"
	"os"
	"testing"

	db "github.com/freedommmoto/test_simplebank/db/sqlc"
	tool "github.com/freedommmoto/test_simplebank/tool"
	_ "github.com/lib/pq"
)

var testQueries *db.Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	config, err := tool.LoadConfig("../")
	//fmt.Printf("%v", config)
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	//error nail pointer here because add code like this testDB, err := sql not  testDB, err = sql
	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testQueries = db.New(testDB)

	os.Exit(m.Run())
}
