package server

import (
	"ckb-net-monitor-log-analyzer/handlers"
	"database/sql"
	"encoding/json"
	"log"
	"time"

	"github.com/lib/pq"
)

// Client construct
type Client struct {
	pgConn *sql.DB
}

// NewClient get client
func NewClient(pgConn *sql.DB) *Client {
	return &Client{pgConn: pgConn}
}

// BulkImport infos to db
func (client *Client) BulkImport(tableName string, infos []handlers.AnalysisInfo, columns ...string) error {
	log.Println("-- Begin Bulk Import --")
	db := client.pgConn
	txn, err := db.Begin()
	if err != nil {
		return err
	}
	notNullColumns := []string{"created_at", "updated_at"}
	columns = append(notNullColumns, columns...)
	stmt, err := txn.Prepare(pq.CopyIn(tableName, columns...))
	if err != nil {
		return err
	}

	for _, info := range infos {
		timestamp := parseTimestamp(info.Timestamp)
		// target_hash, created_at_unixtimestamp,
		jsonDurations, _ := json.Marshal(info.Durations)
		_, err := stmt.Exec(time.Now(), time.Now(), info.TargetHash, timestamp, jsonDurations)
		if err != nil {
			return err
		}
	}
	_, err = stmt.Exec()
	if err != nil {
		return err
	}
	err = stmt.Close()
	if err != nil {
		return err
	}
	err = txn.Commit()
	if err != nil {
		return err
	}
	log.Println("-- End Bulk Import --")
	return nil
}

func parseTimestamp(timestamp int64) (parsedTimestamp int64) {
	secondTimestamp := timestamp / 1000
	parsedTime := time.Unix(secondTimestamp, 0)
	secondsEastOfUTC := int((8 * time.Hour).Seconds())
	beijing := time.FixedZone("Beijing Time", secondsEastOfUTC)
	parsedTimestamp = time.Date(parsedTime.Year(), parsedTime.Month(), parsedTime.Day(), 0, 0, 0, 0, beijing).Unix()
	return
}
