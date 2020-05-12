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

// NewLogAnalyzeService get LogAnalyzeService
func NewLogAnalyzeService(targetLineKeyWord string, pgClient *server.Client) *LogAnalyzeService {
	return &LogAnalyzeService{TargetLineKeyWord: targetLineKeyWord, PGClient: pgClient}
}

// AnalyzeLog can analyze block or transaction propagation delay
func (service *LogAnalyzeService) AnalyzeLog(filePath string, handle func(string, string, map[string]handlers.AnalysisInfo)) error {
	processCount := 1000
	file, err := os.Open(filePath)
	results := make(map[string]handlers.AnalysisInfo)
	defer file.Close()
	if err != nil {
		return err
	}
	buf := bufio.NewReader(file)

	for {
		line, _, err := buf.ReadLine()
		strLine := strings.TrimSpace(string(line))
		handle(strLine, service.TargetLineKeyWord, results)
		saveDataToDB(service, processCount, results)
		if err != nil {
			if err == io.EOF {
				if len(results) > 0 {
					saveDataToDB(service, 1, results)
				}
				return nil
			}
			return err
		}
	}
}

func saveDataToDB(service *LogAnalyzeService, processCount int, results map[string]handlers.AnalysisInfo) {
	infoCompleted := filter(results, func(info handlers.AnalysisInfo) bool {
		// check if 90% duration is calculated
		return info.Durations[17] != 0
	})
	if len(infoCompleted) >= processCount {
		for _, info := range infoCompleted {
			delete(results, info.TargetHash)
		}
		err := service.PGClient.BulkImport("block_propagation_delays", infoCompleted, "block_hash", "created_at_unixtimestamp", "durations")
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
