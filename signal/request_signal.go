package controlsignal

import (
	"fmt"
	"net"
	"strconv"
)

type ControlSignal struct {
	Command   string    `json:"command"`
	FiveTuple FiveTuple `json:"fiveTuple"`
}

func (f ControlSignal) String() string {
	return fmt.Sprintf("command: %s, fiveTuple: %s", f.Command, f.FiveTuple)
}

type FiveTuple struct {
	SrcAddress  string `json:"srcAddress"`
	DestAddress string `json:"destAddress"`
	DestPort    int    `json:"destPort"`
	Protocol    string `json:"protocol"`
}

func (f *FiveTuple) DestJoinedAddress() string {
	return net.JoinHostPort(f.DestAddress, strconv.Itoa(f.DestPort))
}

func (f FiveTuple) String() string {
	return fmt.Sprintf("srcAddress: %s, destAddress: %s, destPort: %d, protocol: %s", f.SrcAddress, f.DestAddress, f.DestPort, f.Protocol)
}
