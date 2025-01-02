package indexer

import (
	"database/sql"

	"github.com/naman1402/geth-indexer/subsrciber"
)

func Index(eventCh chan *subsrciber.Event, db *sql.DB, quit chan bool) {}
