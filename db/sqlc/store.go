package db

import "database/sql"

//
type Srore struct {
	*Queries
	db *sql.DB
}
