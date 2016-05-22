package db

import (
	"database/sql"

	_ "github.com/lib/pq"
)

// Share a connection to the DB between packages
var Conn *sql.DB
