package main

import (
	"flag"
	"log/slog"
	"macl/agent"
	"macl/controller"
	"os"
)

func main() {
	toolType := flag.String("type", "agent", "The type of role to assign to this process. It can be either 'controller' or 'agent'.")
	controlPort := flag.Int("controlPort", 10000, "The port to use for the control signal.")
	profile := flag.String("profile", "", "It defines the profile to use. It will be use to choose config file.")
	debug := flag.Bool("debug", false, "The debug mode. It will be used to set the log level. true for debug mode, false for info mode.")
	flag.Parse()

	// Initialize the logger
	var level slog.Level
	if *debug {
		level = slog.LevelDebug
	} else {
		level = slog.LevelInfo
	}
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	logger := slog.New(handler)

	switch *toolType {
	case "agent":
		agent := agent.NewAgent(logger, *controlPort)
		agent.Start()
	case "controller":
		fileTuples := controller.NewControllerConfig(logger).LoadConfig(*profile)
		controller := controller.NewController(logger, *controlPort)
		controller.Start(fileTuples)
	default:
	}

}
