package data

import (
	"database/sql"
	"fmt"
	_ "github.com/microsoft/go-mssqldb"
	"log"
)

var db *sql.DB
var server = "bs054.database.windows.net"
var port = 1433
var user = "mehedi"
var password = "database@123"
var database = "test"

func main() {
	connectionString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
		server, user, password, port, database)

	var err error
	db, err := sql.Open("sqlserver", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	fmt.Println("Connected")
}
