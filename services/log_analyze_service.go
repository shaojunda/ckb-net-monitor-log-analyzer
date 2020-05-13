package services

import (
	"bufio"
	"ckb-net-monitor-log-analyzer/handlers"
	"ckb-net-monitor-log-analyzer/server"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

// LogAnalyzeService construct
type LogAnalyzeService struct {
	TargetLineKeyWord string
	ProcessCount      int
	PGClient          *server.Client
}

type dbTableInfo struct {
	tableName, columnName string
}

type processResult struct {
	FileName string
	Position int64
	Results  map[string]handlers.AnalysisInfo
}

// NewLogAnalyzeService get LogAnalyzeService
func NewLogAnalyzeService(targetLineKeyWord string, processCount int, pgClient *server.Client) *LogAnalyzeService {
	return &LogAnalyzeService{TargetLineKeyWord: targetLineKeyWord, ProcessCount: processCount, PGClient: pgClient}
}

// AnalyzeLog can analyze block or transaction propagation delay
func (service *LogAnalyzeService) AnalyzeLog(filePath string, handle func(string, string, map[string]handlers.AnalysisInfo)) error {
	var tableInfo dbTableInfo
	if service.TargetLineKeyWord == "compact_block:" {
		tableInfo = dbTableInfo{tableName: "block_propagation_delays", columnName: "block_hash"}
	} else {
		tableInfo = dbTableInfo{tableName: "transaction_propagation_delays", columnName: "tx_hash"}
	}
	start, results := initProcessInfo(tableInfo)

	return readFileWithScanner(filePath, start, service.ProcessCount, service, handle, results, tableInfo)
}

func setupCloseHandler(processInfo *processResult) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\r Get Interrupt Signal")
		fmt.Println("Interrupt position: ", processInfo.Position)
		saveProcessInfo(processInfo)
		os.Exit(0)
	}()
}

func initProcessInfo(tableInfo dbTableInfo) (start int64, results map[string]handlers.AnalysisInfo) {
	fileName := tableInfo.tableName + ".json"
	processIno := processResult{FileName: fileName}
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		return 0, make(map[string]handlers.AnalysisInfo)
	}

	json.Unmarshal(file, &processIno)
	return processIno.Position, processIno.Results
}

func saveProcessInfo(processIno *processResult) {
	jsonString, err := json.Marshal(processIno)
	if err != nil {
		log.Println("error: ", err)
	}
	_ = ioutil.WriteFile(processIno.FileName, jsonString, 0644)
}

func readFileWithScanner(filePath string, start int64, processCount int, service *LogAnalyzeService, handle func(string, string, map[string]handlers.AnalysisInfo), results map[string]handlers.AnalysisInfo, tableInfo dbTableInfo) error {
	log.Printf("--%s SCANNER, start: %d\n", tableInfo.tableName, start)
	jsonFileName := tableInfo.tableName + ".json"
	file, err := os.Open(filePath)
	pos := start
	processIno := processResult{FileName: jsonFileName, Position: pos, Results: results}
	setupCloseHandler(&processIno)
	defer func() {
		file.Close()
		saveProcessInfo(&processIno)
		log.Printf("%s Done.\n", tableInfo.tableName)
	}()
	if err != nil {
		return err
	}
	if _, err = file.Seek(start, 0); err != nil {
		return err
	}

	scanner := bufio.NewScanner(file)

	scanLines := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		advance, token, err = bufio.ScanLines(data, atEOF)
		pos += int64(advance)
		processIno.Position = pos
		return
	}
	scanner.Split(scanLines)

	for scanner.Scan() {
		strLine := strings.TrimSpace(scanner.Text())
		handle(strLine, service.TargetLineKeyWord, results)
		saveDataToDB(service, processCount, tableInfo, results)
	}

	processCount = 1
	if len(results) > 0 {
		saveDataToDB(service, processCount, tableInfo, results)
	}

	return scanner.Err()
}

func readFile(filePath string, processCount int, service *LogAnalyzeService, handle func(string, string, map[string]handlers.AnalysisInfo), results map[string]handlers.AnalysisInfo, tableInfo dbTableInfo) error {
	file, err := os.Open(filePath)
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
		err := service.PGClient.BulkImport(tableInfo.tableName, infoCompleted, tableInfo.columnName, "created_at_unixtimestamp", "durations")
		if err != nil {
			log.Println(err.Error())
		} else {
			for _, info := range infoCompleted {
				delete(results, info.TargetHash)
			}
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
