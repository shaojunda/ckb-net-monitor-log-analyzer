package services

import (
	"bufio"
	"ckb-net-monitor-log-analyzer/handlers"
	"ckb-net-monitor-log-analyzer/server"
	"fmt"
	"io"
	"os"
	"strings"
)

// LogAnalyzeService construct
type LogAnalyzeService struct {
	TargetLineKeyWord string
	PGClient          *server.Client
}

type dbTableInfo struct {
	tableName, columnName string
}

// NewLogAnalyzeService get LogAnalyzeService
func NewLogAnalyzeService(targetLineKeyWord string, pgClient *server.Client) *LogAnalyzeService {
	return &LogAnalyzeService{TargetLineKeyWord: targetLineKeyWord, PGClient: pgClient}
}

// AnalyzeLog can analyze block or transaction propagation delay
func (service *LogAnalyzeService) AnalyzeLog(filePath string, handle func(string, string, map[string]handlers.AnalysisInfo)) error {
	var tableInfo dbTableInfo
	processCount := 1000
	file, err := os.Open(filePath)
	results := make(map[string]handlers.AnalysisInfo)
	if service.TargetLineKeyWord == "compact_block:" {
		tableInfo = dbTableInfo{tableName: "block_propagation_delays", columnName: "block_hash"}
	} else {
		tableInfo = dbTableInfo{tableName: "transaction_propagation_delays", columnName: "tx_hash"}
	}
	defer file.Close()
	if err != nil {
		return err
	}
	buf := bufio.NewReader(file)

	for {
		line, _, err := buf.ReadLine()
		strLine := strings.TrimSpace(string(line))
		handle(strLine, service.TargetLineKeyWord, results)
		saveDataToDB(service, processCount, tableInfo, results)
		if err != nil {
			if err == io.EOF {
				processCount = 1
				if len(results) > 0 {
					saveDataToDB(service, processCount, tableInfo, results)
				}
				return nil
			}
			return err
		}
	}
}

func saveDataToDB(service *LogAnalyzeService, processCount int, tableInfo dbTableInfo, results map[string]handlers.AnalysisInfo) {
	infoCompleted := filter(results, func(info handlers.AnalysisInfo) bool {
		// check if 90% duration is calculated
		return info.Durations[17] != 0
	})
	if len(infoCompleted) >= processCount {
		for _, info := range infoCompleted {
			delete(results, info.TargetHash)
		}
		err := service.PGClient.BulkImport(tableInfo.tableName, infoCompleted, tableInfo.columnName, "created_at_unixtimestamp", "durations")
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

func filter(analysisInfo map[string]handlers.AnalysisInfo, f func(handlers.AnalysisInfo) bool) []handlers.AnalysisInfo {
	infos := make([]handlers.AnalysisInfo, 0, 0)
	for _, value := range analysisInfo {
		if f(value) {
			infos = append(infos, value)
		}
	}

	return infos
}
