package services

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

const targetLineLength = 11

// LogAnalyzeService construct
type LogAnalyzeService struct {
	TargetLineKeyWord string
}

func NewLogAnalyzeService(targetLineKeyWord string) *LogAnalyzeService {
	return &LogAnalyzeService{TargetLineKeyWord: targetLineKeyWord}
}

type analystInfo struct {
	Count     int
	Timestamp int64
	Durations [18]int64
}

func parsePeers(peers string) (peersInt int) {
	peersInt, err := strconv.Atoi(peers)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
	}
	return
}

func parseTimestamp(timestamp string) (unixTimestamp int64) {
	parsedTime, err := time.Parse("2006-01-02 15:04:05.000 -0700 MST", timestamp)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
	}
	unixTimestamp = parsedTime.UnixNano() / int64(time.Millisecond)
	return
}

func calculateDuration(analystInfo *analystInfo, peers int, unixTimestamp int64) {
	pairs := map[int]int{5: 0, 10: 1, 15: 2, 20: 3, 25: 4, 30: 5, 35: 6, 40: 7, 45: 8, 50: 9, 55: 10, 60: 11, 65: 12, 70: 13, 75: 14, 80: 15, 85: 16, 90: 17}
	for key, value := range pairs {
		if analystInfo.Count*100 >= peers*key && analystInfo.Durations[value] == 0 {
			duration := unixTimestamp - analystInfo.Timestamp
			analystInfo.Durations[value] = duration
		}
	}
}

func (service *LogAnalyzeService) Handle(line string, results map[string]analystInfo) {
	targetLine := strings.Split(line, " ")
	if len(targetLine) == targetLineLength && targetLine[7] == &service.TargetLineKeyWord {
		targetHash := strings.TrimRight(targetLine[8], ",")
		peers := parsePeers(targetLine[10])
		timestamp := targetLine[0] + " " + targetLine[1] + " " + strings.Replace(targetLine[2], ":", "", 1) + " CST"
		unixTimestamp := parseTimestamp(timestamp)
		if val, ok := results[targetHash]; ok {
			val.Count++
			calculateDuration(&val, peers, unixTimestamp)
			results[targetHash] = val
		} else {
			results[targetHash] = analystInfo{Count: 1, Timestamp: unixTimestamp}
		}
	}
}

func parseLog(filePath string, handle func(string, map[string]analystInfo)) error {
	file, err := os.Open(filePath)
	results := make(map[string]analystInfo)
	defer file.Close()
	if err != nil {
		return err
	}
	buf := bufio.NewReader(file)

	for {
		line, _, err := buf.ReadLine()
		strLine := strings.TrimSpace(string(line))
		handle(strLine, results)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
	}
}
