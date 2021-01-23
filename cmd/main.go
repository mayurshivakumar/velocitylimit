package main

import (
	"bufio"
	"encoding/json"
	"os"

	"velocitylimit/models"

	"velocitylimit/cache"
	"velocitylimit/service"

	"velocitylimit/config"

	"github.com/sirupsen/logrus"
)

func main() {
	cfg := config.ParseConfig()
	cache := cache.NewCache()
	inputFile := OpenFile(cfg)
	outputFile := CreateFile(cfg)

	defer inputFile.Close()
	defer outputFile.Close()

	scanner := bufio.NewScanner(inputFile)
	writer := bufio.NewWriter(outputFile)
	for scanner.Scan() {
		// get request
		request, err := models.NewRequest(scanner.Text())
		if err != nil {
			logrus.Errorf("Error reading file:%v", err)
			continue
		}
		// attempt to load
		response := service.AttemptLoad(request, cfg, cache)
		resBytes, err := json.Marshal(response)
		if err != nil {
			logrus.Errorf("Error marshalling json:%v", err)
			continue
		}
		// write to file
		if _, err = writer.WriteString(string(resBytes) + "\n"); err != nil {
			logrus.Panicf("Error writing to file file:%v", err)
			continue
		}

		// error reading file
		if err := scanner.Err(); err != nil {
			logrus.Panicf("Error reading file:%v", err)
		}
	}
	writer.Flush()
}

func OpenFile(config *config.Configurations) *os.File {
	// TODO: Path for the file needs to be handled better
	input, err := os.Open("../" + config.VelocityLimit.InputFile)
	if err != nil {
		logrus.Panicf("Unable to open file: %s", err)
	}
	return input
}

func CreateFile(config *config.Configurations) *os.File {
	// TODO: Path for the file needs to be handled better
	output, err := os.Create("../" + config.VelocityLimit.OutputFile)
	if err != nil {
		logrus.Panicf("Unable to open file: %s", err)
	}
	return output
}
