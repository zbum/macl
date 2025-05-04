package main

import (
	"flag"
	"macl/agent"
)

func main() {
	toolType := flag.String("type", "agent", "The type of role to assign to this process. It can be either 'controller' or 'agent'.")
	controlPort := flag.Int("controlPort", 10000, "The port to use for the control signal.")
	//protocol := flag.String("protocol", "tcp", "The protocol assign to this process. It can be either 'tcp' or 'udp'.")
	//profile := flag.String("profile", "", "It defines the profile to use. It will be use to choose config file.")
	flag.Parse()

	switch *toolType {
	case "agent":
		agent := agent.NewAgent(*controlPort)
		agent.Start()
	case "controller":

	default:
	}

}
