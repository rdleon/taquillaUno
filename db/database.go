package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

// Share a connection to the DB between packages
var Conn *sql.DB

// Connects to a test database in a local postgres instance
func InitTestDB() *sql.DB {
	// TODO: Read test info from config file
	var err error

	var conf = map[string]string{
		"db_user":   "taquilla",
		"db_passwd": "secret",
		"db_name":   "taquillatest",
		"db_host":   "localhost",
		"listen":    "127.0.0.1:8080",
		"secret":    "secretkeyVerySecret",
	}

	dbConf := fmt.Sprintf(
		"user=%s dbname=%s password=%s",
		conf["db_user"], conf["db_name"], conf["db_passwd"],
	)

	Conn, err = sql.Open("postgres", dbConf)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error connecting to the database", err)
		return nil
	}

	if err != nil {
		fmt.Println(err)
	}

	return Conn
}

// Reads the tables from the SQL files and creates them in the
// test database.
func CreateTestTables() (err error) {
	err = createTable("users")
	return
}

// Delete the relations on the test database
func ClearTestTables() (err error) {
	err = dropTable("users")
	return
}

func createTable(table string) (err error) {
	var b []byte
	var count int

	fh, err := os.Open("db/" + table + ".sql")

	if err != nil {
		//TODO: Check the error and return "Table SQL file (<table>.sql) not found
		fmt.Println(err)
		return
	}

	defer fh.Close()

	b = make([]byte, 1024)
	count, err = fh.Read(b)

	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = Conn.Query(string(b[:count]))
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Table " + table + " created")

	return
}

func dropTable(table string) (err error) {
	_, err = Conn.Query("DROP TABLE IF EXISTS " + table)
	return
}

func addDummyData() (err error) {
}
