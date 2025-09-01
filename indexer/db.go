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

	// Create transfer table if it doesn't exist
	if err := createTransferTable(db); err != nil {
		log.Printf("Warning: failed to create transfer table: %v", err)
	}

	return db, nil
}

func createTransferTable(db *sql.DB) error {
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS transfer (
		id SERIAL PRIMARY KEY,
		"name" VARCHAR(50) NOT NULL,
		"blockNumber" BIGINT NOT NULL,
		"txnHash" VARCHAR(66) NOT NULL,
		"contract" VARCHAR(42) NOT NULL,
		"from" VARCHAR(42) NOT NULL,
		"to" VARCHAR(42) NOT NULL,
		"value" NUMERIC NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		UNIQUE("txnHash", "contract", "from", "to", "value")
	);`

	_, err := db.Exec(createTableQuery)
	if err != nil {
		return fmt.Errorf("failed to create transfer table: %v", err)
	}

	// Create indexes for better performance
	indexes := []string{
		`CREATE INDEX IF NOT EXISTS idx_transfer_contract ON transfer("contract")`,
		`CREATE INDEX IF NOT EXISTS idx_transfer_from ON transfer("from")`,
		`CREATE INDEX IF NOT EXISTS idx_transfer_to ON transfer("to")`,
		`CREATE INDEX IF NOT EXISTS idx_transfer_block ON transfer("blockNumber")`,
		`CREATE INDEX IF NOT EXISTS idx_transfer_txn ON transfer("txnHash")`,
	}

	for _, indexQuery := range indexes {
		if _, err := db.Exec(indexQuery); err != nil {
			log.Printf("Warning: failed to create index: %v", err)
		}
	}

	log.Println("Transfer table and indexes created successfully")
	return nil
}

// executeQuery executes a parameterized query with the provided args.
func executeQuery(db *sql.DB, query string, args ...interface{}) {
	_, err := db.Exec(query, args...)
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
