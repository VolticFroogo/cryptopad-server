package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/VolticFroogo/config"
	"github.com/gchaincl/dotsql"
	_ "github.com/go-sql-driver/mysql" // MySQL Driver.
)

var (
	// SQL is the global database DB connection.
	SQL *sql.DB

	// Dot is all of the loaded queries.
	Dot *dotsql.DotSql
)

// Config is the config structure.
type Config struct {
	Name, Password, Protocol, Location, Database, QueriesDirectory string
}

// Init initialises the database.
func Init(configDirectory string) (err error) {
	// Load the config.
	cfg := Config{}
	err = config.Load(configDirectory, &cfg)
	if err != nil {
		return
	}

	// Log that we are connecting to the database.
	log.Print("Connecting to database.")

	// Create the connection string from the config.
	connection := fmt.Sprintf("%v:%v@%v(%v)/%v", cfg.Name, cfg.Password, cfg.Protocol, cfg.Location, cfg.Database)

	// Open the SQL connection.
	SQL, err = sql.Open("mysql", connection)
	if err != nil {
		return
	}

	// Load the queries SQL file.
	Dot, err = dotsql.LoadFromFile(cfg.QueriesDirectory)
	return
}
