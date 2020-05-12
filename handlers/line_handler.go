package handlers

import (
	"log"
	"strconv"
	"strings"
	"time"
)

const targetLineLength = 11

// AnalysisInfo construct
type AnalysisInfo struct {
	Count      int
	Timestamp  int64
	Durations  [18]int64
	TargetHash string
}

// Handle function process each line of data in the log
func Handle(line string, keyword string, results map[string]AnalysisInfo) {
	targetLine := strings.Split(line, " ")
	if len(targetLine) == targetLineLength && targetLine[7] == keyword {
		targetHash := strings.TrimRight(targetLine[8], ",")
		peers := parsePeers(targetLine[10])
		timestamp := targetLine[0] + " " + targetLine[1] + " " + strings.Replace(targetLine[2], ":", "", 1) + " CST"
		unixTimestamp := parseTimestamp(timestamp)
		if val, ok := results[targetHash]; ok {
			val.Count++
			calculateDuration(&val, peers, unixTimestamp)
			results[targetHash] = val
		} else {
			results[targetHash] = AnalysisInfo{Count: 1, Timestamp: unixTimestamp, TargetHash: targetHash}
		}
	}
}

func parsePeers(peers string) (peersInt int) {
	peersInt, err := strconv.Atoi(peers)
	if err != nil {
		log.Fatal("Parse Peers To Int Failed", err)
	}
	return
}

func parseTimestamp(timestamp string) (unixTimestamp int64) {
	parsedTime, err := time.Parse("2006-01-02 15:04:05.000 -0700 MST", timestamp)
	if err != nil {
		log.Fatal("Parse Timestamp Failed", err)
	}
	unixTimestamp = parsedTime.UnixNano() / int64(time.Millisecond)
	return
}

func calculateDuration(analysisInfo *AnalysisInfo, peers int, unixTimestamp int64) {
	// key is percentage, value is analysisInfo.Durations index
	pairs := map[int]int{5: 0, 10: 1, 15: 2, 20: 3, 25: 4, 30: 5, 35: 6, 40: 7, 45: 8, 50: 9, 55: 10, 60: 11, 65: 12, 70: 13, 75: 14, 80: 15, 85: 16, 90: 17}
	for key, value := range pairs {
		if analysisInfo.Count*100 >= peers*key && analysisInfo.Durations[value] == 0 {
			duration := unixTimestamp - analysisInfo.Timestamp
			analysisInfo.Durations[value] = duration
		}
	}
}
