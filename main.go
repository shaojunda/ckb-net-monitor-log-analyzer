package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"

	"ckb-net-monitor-log-analyzer/handlers"
	"ckb-net-monitor-log-analyzer/server"
	"ckb-net-monitor-log-analyzer/services"

	_ "github.com/lib/pq"
	"gopkg.in/yaml.v3"
)

func main() {
	var c config
	config := c.getConfig()
	pgConn, err := getConn(config.PgPort, config.PgHost, config.PgUser, config.PgPassword, config.PgDBName)
	if err != nil {
		panic(err)
	}
	client := server.NewClient(pgConn)
	// blockKeyWord := "compact_block:"
	// blockAnalyzeService := services.NewLogAnalyzeService(blockKeyWord, client)
	// blockAnalyzeService.AnalyzeLog(config.MonitorLogFilePath, handlers.Handle)
	transactionKeyWord := "relay_transaction_hashes:"
	transactionAnalyzeService := services.NewLogAnalyzeService(transactionKeyWord, client)
	transactionAnalyzeService.AnalyzeLog(config.MonitorLogFilePath, handlers.Handle)
}

type config struct {
	MonitorLogFilePath string `yaml:"monitor_log_file_path"`
	PgHost             string `yaml:"pg_host"`
	PgPort             int    `yaml:"pg_port"`
	PgUser             string `yaml:"pg_user"`
	PgPassword         string `yaml:"pg_password"`
	PgDBName           string `yaml:"pg_db_name"`
}

func (c *config) getConfig() *config {
	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Fatal("Read Config File Failed", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatal("Parse Config File Failed", err)
	}

	return c
}

func getConn(pgPort int, pgHost, pgUser, pgPassword, pgDBName string) (*sql.DB, error) {
	log.Println("Connecting PostgreSQL...")
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", pgHost, pgPort, pgUser, pgPassword, pgDBName)
	log.Println("connStr: ", connStr)
	db, err := sql.Open("postgres", connStr)
	// limit the number of idle connections in the pool
	db.SetMaxIdleConns(100)
	// limit the number of total open connections to the database.
	db.SetMaxOpenConns(200)
	if err != nil {
		log.Fatal("Connect PG Failed", err)
		return nil, err
	}
	return db, nil
}
