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

func generateQuery(table string, param *subsrciber.Event) string {

	columns := len(param.Data)
	fieldSlice := make([]string, 0, columns)
	fields := "name, blockNumber, blockHash, contract, "

	for fields := range param.Data {
		fieldSlice = append(fieldSlice, fields)
	}

	fields += strings.Join(fieldSlice, " ,")

	//
	values := addValues(fieldSlice, param)
	return fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s)`, table, fields, values)

}

func addValues(keys []string, params *subsrciber.Event) string {
	values := make([]string, 0, len(keys)+4)
	values = append(values, fmt.Sprintf("%v", params.Name), fmt.Sprintf("%v", params.BlockNumber), fmt.Sprintf("%v", params.BlockHash), fmt.Sprintf("%v", params.Contract))

	for _, k := range keys {
		values = append(values, fmt.Sprintf("%v", params.Data[k]))
	}

	return strings.Join(values, " ,")
}
