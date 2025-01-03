package indexer

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/naman1402/geth-indexer/subsrciber"
)

func Index(eventCh chan *subsrciber.Event, db *sql.DB, quit chan bool) {

	for {
		select {
		case e := <-eventCh:
			query := generateQuery(strings.ToLower(e.Name), e)
			go executeQuery(db, query)
		case q := <-quit:
			if q {
				return
			}
		}
	}
}

// generateQuery constructs an SQL INSERT statement for the given table and event parameters.
// It generates the column names and values based on the event data, and returns the complete SQL query.
func generateQuery(table string, param *subsrciber.Event) string {

	columns := len(param.Data)
	// used to store  the names of the fields in the event data
	fieldSlice := make([]string, 0, columns)
	// primary fields, represents basic information about the event that will be included in the SQL query
	fields := "name, blockNumber, blockHash, contract, "

	// iterate through the params.Data and append to the fieldSlice
	// results in a complete list of all the fields that will be included in the SQL query
	for fields := range param.Data {
		fieldSlice = append(fieldSlice, fields)
	}
	// concatenate the fields and fieldSlice into a single string
	fields += strings.Join(fieldSlice, " ,")

	// constructs a string of values corresponding to the fields, which will be used in the VALUES clause of the SQL query.
	// Unlike fieldSlice, this will contain the actual values from the event data.
	values := addValues(fieldSlice, param)
	return fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s)`, table, fields, values)

}

// addValues constructs a comma-separated string of values from the given keys and event parameters.
// The resulting string can be used in an SQL INSERT statement.
func addValues(keys []string, params *subsrciber.Event) string {

	values := make([]string, 0, len(keys)+4)
	// Appending 4 default values to the slice
	values = append(values, fmt.Sprintf("%v", params.Name), fmt.Sprintf("%v", params.BlockNumber), fmt.Sprintf("%v", params.BlockHash), fmt.Sprintf("%v", params.Contract))
	// corresponding value from the params.Data map and appends its formatted string representation to the values slice.
	for _, k := range keys {
		values = append(values, fmt.Sprintf("%v", params.Data[k]))
	}

	return strings.Join(values, " ,")
}
