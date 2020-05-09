package services

import (
	"bufio"
	"ckb-net-monitor-log-analyzer/handlers"
	"io"
	"os"
	"strings"
)

// LogAnalyzeService construct
type LogAnalyzeService struct {
	TargetLineKeyWord string
}

// NewLogAnalyzeService get LogAnalyzeService
func NewLogAnalyzeService(targetLineKeyWord string) *LogAnalyzeService {
	return &LogAnalyzeService{TargetLineKeyWord: targetLineKeyWord}
}

// AnalyzeLog can analyze block or transaction propagation delay
func (service *LogAnalyzeService) AnalyzeLog(filePath string, handle func(string, string, map[string]handlers.AnalysisInfo)) error {
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
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
	}
}
