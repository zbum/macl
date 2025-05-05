package controller

import (
	"encoding/csv"
	"fmt"
	"io"
	"log/slog"
	controlsignal "macl/signal"
	"os"
	"strconv"
)

type ControllerConfig struct {
	log *slog.Logger
}

func NewControllerConfig(log *slog.Logger) *ControllerConfig {
	return &ControllerConfig{
		log: log,
	}
}

func (c *ControllerConfig) LoadConfig(profile string) []*controlsignal.FiveTuple {
	log := c.log

	var configFile *os.File
	var err error
	if profile == "" {
		configFile, err = os.Open("config.csv")
	} else {
		configFile, err = os.Open(fmt.Sprintf("config-%s.csv", profile))
	}
	if err != nil {
		log.Error("[macl-controller] config file failed", "error", err)
		return nil
	}
	defer configFile.Close()

	r := csv.NewReader(configFile)
	var fiveTuples []*controlsignal.FiveTuple
	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}

		if len(row) < 4 {
			log.Error("[macl-controller] invalid row", "row", row)
			continue
		}

		destPort, err := strconv.Atoi(row[3])
		if err != nil {
			log.Error("[macl-controller] invalid row", "row", row, "error", err)
			continue
		}

		fiveTuples = append(fiveTuples, &controlsignal.FiveTuple{
			TxId:        row[0],
			SrcAddress:  row[1],
			DestAddress: row[2],
			DestPort:    destPort,
			Protocol:    row[4],
		})
	}

	return fiveTuples
}
