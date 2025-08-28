package indexer

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"

	"github.com/naman1402/geth-indexer/cli"
)

func Connect(options cli.DatabaseConfig) (*sql.DB, error) {

	postgreSQLInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", options.DBHost, options.DBPort, options.DBUser, options.DBPassword, options.DBName)
	db, err := sql.Open("postgres", postgreSQLInfo)
	if err != nil {
		log.Printf("failed to connect to database: %s", err)
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		log.Printf("Failed to ping database: %s", err)
		return nil, err
	}
	fmt.Printf("Ping successful: connected to the database %s and port %d\n", options.DBName, options.DBPort)
	return db, nil
}

func executeQuery(db *sql.DB, query string) {
	_, err := db.Exec(query)
	if err != nil {
		log.Println("failed to execute query:", query, err)
	}
}

// func indexingCheck(db *sql.DB, relation, column string) {

// 	index := relation + "_" + column + "_idx"
// 	_, err := db.Exec("CREATE INDEX %s IF NOT EXISTS  %s ON %s (%s)", index, index, relation, column)
// 	if err != nil {
// 		log.Println(err)
// 	}
// }
