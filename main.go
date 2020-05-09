package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

const targetLineLength = 11
const targetLineKeyWord = "compact_block:"

func main() {
	var c config
	config := c.getConfig()
	parseLog(config.MonitorLogFilePath, handle)
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

func handle(line string, results map[string]analystInfo) {
	targetLine := strings.Split(line, " ")
	if len(targetLine) == targetLineLength && targetLine[7] == targetLineKeyWord {
		blockHash := strings.TrimRight(targetLine[8], ",")
		peers := parsePeers(targetLine[10])
		timestamp := targetLine[0] + " " + targetLine[1] + " " + strings.Replace(targetLine[2], ":", "", 1) + " CST"
		unixTimestamp := parseTimestamp(timestamp)
		if val, ok := results[blockHash]; ok {
			val.Count++
			if val.Count*100 >= peers*5 && val.Durations[0] == 0 {
				duration := unixTimestamp - val.Timestamp
				val.Durations[0] = duration
				fmt.Printf("5%% blockHash: %s duration: %d\n", blockHash, duration)
			}
			if val.Count*100 >= peers*10 && val.Durations[1] == 0 {
				duration := unixTimestamp - val.Timestamp
				val.Durations[1] = duration
				fmt.Printf("10%% blockHash: %s duration: %d\n", blockHash, duration)
			}
			if val.Count*100 >= peers*15 && val.Durations[2] == 0 {
				duration := unixTimestamp - val.Timestamp
				val.Durations[2] = duration
				fmt.Printf("15%% blockHash: %s duration: %d\n", blockHash, duration)
			}
			if val.Count*100 >= peers*20 && val.Durations[3] == 0 {
				duration := unixTimestamp - val.Timestamp
				val.Durations[3] = duration
				fmt.Printf("20%% blockHash: %s duration: %d\n", blockHash, duration)
			}
			if val.Count*100 >= peers*25 && val.Durations[4] == 0 {
				duration := unixTimestamp - val.Timestamp
				val.Durations[4] = duration
				fmt.Printf("25%% blockHash: %s duration: %d\n", blockHash, duration)
			}
			if val.Count*100 >= peers*30 && val.Durations[5] == 0 {
				duration := unixTimestamp - val.Timestamp
				val.Durations[5] = duration
				fmt.Printf("30%% blockHash: %s duration: %d\n", blockHash, duration)
			}
			if val.Count*100 >= peers*35 && val.Durations[6] == 0 {
				duration := unixTimestamp - val.Timestamp
				val.Durations[6] = duration
				fmt.Printf("35%% blockHash: %s duration: %d\n", blockHash, duration)
			}
			if val.Count*100 >= peers*40 && val.Durations[7] == 0 {
				duration := unixTimestamp - val.Timestamp
				val.Durations[7] = duration
				fmt.Printf("40%% blockHash: %s duration: %d\n", blockHash, duration)
			}
			if val.Count*100 >= peers*45 && val.Durations[8] == 0 {
				duration := unixTimestamp - val.Timestamp
				val.Durations[8] = duration
				fmt.Printf("45%% blockHash: %s duration: %d\n", blockHash, duration)
			}
			if val.Count*100 >= peers*50 && val.Durations[9] == 0 {
				duration := unixTimestamp - val.Timestamp
				val.Durations[9] = duration
				fmt.Printf("50%% blockHash: %s duration: %d\n", blockHash, duration)
			}
			if val.Count*100 >= peers*55 && val.Durations[10] == 0 {
				duration := unixTimestamp - val.Timestamp
				val.Durations[10] = duration
				fmt.Printf("55%% blockHash: %s duration: %d\n", blockHash, duration)
			}
			if val.Count*100 >= peers*60 && val.Durations[11] == 0 {
				duration := unixTimestamp - val.Timestamp
				val.Durations[11] = duration
				fmt.Printf("60%% blockHash: %s duration: %d\n", blockHash, duration)
			}
			if val.Count*100 >= peers*65 && val.Durations[12] == 0 {
				duration := unixTimestamp - val.Timestamp
				val.Durations[12] = duration
				fmt.Printf("65%% blockHash: %s duration: %d\n", blockHash, duration)
			}
			if val.Count*100 >= peers*70 && val.Durations[13] == 0 {
				duration := unixTimestamp - val.Timestamp
				val.Durations[13] = duration
				fmt.Printf("70%% blockHash: %s duration: %d\n", blockHash, duration)
			}
			if val.Count*100 >= peers*75 && val.Durations[14] == 0 {
				duration := unixTimestamp - val.Timestamp
				val.Durations[14] = duration
				fmt.Printf("75%% blockHash: %s duration: %d\n", blockHash, duration)
			}
			if val.Count*100 >= peers*80 && val.Durations[15] == 0 {
				duration := unixTimestamp - val.Timestamp
				val.Durations[15] = duration
				fmt.Printf("80%% blockHash: %s duration: %d\n", blockHash, duration)
			}
			if val.Count*100 >= peers*85 && val.Durations[16] == 0 {
				duration := unixTimestamp - val.Timestamp
				val.Durations[16] = duration
				fmt.Printf("85%% blockHash: %s duration: %d\n", blockHash, duration)
			}
			if val.Count*100 >= peers*90 && val.Durations[17] == 0 {
				duration := unixTimestamp - val.Timestamp
				val.Durations[17] = duration
				fmt.Printf("90%% blockHash: %s duration: %d\n", blockHash, duration)
			}
			results[blockHash] = val
		} else {
			results[blockHash] = analystInfo{Count: 1, Timestamp: unixTimestamp}
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
