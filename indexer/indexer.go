package indexer

import (
	"database/sql"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/naman1402/geth-indexer/subsrciber"
)

// Index is the main function that listens for events on the eventCh channel, generates SQL queries
// using the generateQuery function, and executes those queries asynchronously using the executeQuery function.
// The function runs in an infinite loop, waiting for events or a quit signal on the quit channel.
// When a quit signal is received, the function returns.
func Index(eventCh chan *subsrciber.Event, db *sql.DB, quit chan bool) {

	for {
		select {
		case e := <-eventCh:
			query, args := generateQuery(strings.ToLower(e.Name), e)
			// log.Printf("[Index] received event from eventCh, creating query and executing it. Query: %s, Args: %v", query, args)
			log.Printf("[Index] received event from eventCh, creating query and executing it")
			go executeQuery(db, query, args...)
		case q := <-quit:
			if q {
				return
			}
		}
	}
}

// generateQuery constructs an SQL INSERT statement for the given table and event parameters.
// It generates the column names and values based on the event data, and returns the complete SQL query.
func generateQuery(table string, param *subsrciber.Event) (string, []interface{}) {

	// collect field names from event data (iteration order is indeterminate)
	fieldSlice := make([]string, 0, len(param.Data))
	for field := range param.Data {
		fieldSlice = append(fieldSlice, field)
	}

	// base columns
	allCols := []string{"name", "blockNumber", "txnHash", "contract"}
	allCols = append(allCols, fieldSlice...)

	// quoted column list to avoid reserved word collisions
	colsQuoted := make([]string, 0, len(allCols))
	for _, c := range allCols {
		// quote identifiers to allow reserved words like from/to as column names
		colsQuoted = append(colsQuoted, fmt.Sprintf(`"%s"`, c))
	}
	colsStr := strings.Join(colsQuoted, ", ")

	total := len(allCols)
	placeholders := make([]string, 0, total)
	for i := 1; i <= total; i++ {
		placeholders = append(placeholders, fmt.Sprintf("$%d", i))
	}
	phStr := strings.Join(placeholders, ", ")

	// build args in the same order as allCols
	args := make([]interface{}, 0, total)
	args = append(args, param.Name)
	args = append(args, param.BlockNumber)
	args = append(args, fmt.Sprintf("%s", param.TxnHash))
	args = append(args, fmt.Sprintf("%s", param.Contract))
	for _, k := range fieldSlice {
		v := param.Data[k]
		switch val := v.(type) {
		case *big.Int:
			args = append(args, val.String())
		default:
			args = append(args, val)
		}
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) ON CONFLICT (\"txnHash\", \"contract\", \"from\", \"to\", \"value\") DO NOTHING", table, colsStr, phStr)
	return query, args

}
