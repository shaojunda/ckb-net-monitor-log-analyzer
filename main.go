package main

import (
	"fmt"
	"io/ioutil"

	"ckb-net-monitor-log-analyzer/handlers"
	services "ckb-net-monitor-log-analyzer/services"

	"gopkg.in/yaml.v3"
)

func main() {
	var c config
	config := c.getConfig()
	blockKeyWord := "compact_block:"
	// transactionKeyWord := "relay_transaction_hashes:"
	blockAnalyzeService := services.NewLogAnalyzeService(blockKeyWord)
	// transactionAnalyzeService := services.NewLogAnalyzeService(transactionKeyWord)
	blockAnalyzeService.AnalyzeLog(config.MonitorLogFilePath, handlers.Handle)
}

type config struct {
	MonitorLogFilePath string `yaml:"monitor_log_file_path"`
}

func (c *config) getConfig() *config {
	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		fmt.Println(err.Error())
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		fmt.Println(err.Error())
	}

	return c
}
