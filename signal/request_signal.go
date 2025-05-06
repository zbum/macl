package controlsignal

import (
	"fmt"
	netutil "macl/net"
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
	TxId        string `json:"txId"`
	SrcAddress  string `json:"srcAddress"`
	DestAddress string `json:"destAddress"`
	DestPort    int    `json:"destPort"`
	Protocol    string `json:"protocol"`
}

func (f FiveTuple) DestJoinedAddress() string {
	return net.JoinHostPort(f.DestAddress, strconv.Itoa(f.DestPort))
}

func (f FiveTuple) String() string {
	return fmt.Sprintf("txId: %s, srcAddress: %s, destAddress: %s, destPort: %d, protocol: %s", f.TxId, f.SrcAddress, f.DestAddress, f.DestPort, f.Protocol)
}

func (f FiveTuple) AmISource() (bool, error) {
	myHostIp, err := netutil.IsMyActiveHostIp(f.SrcAddress)
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	return myHostIp, nil
}

func (f FiveTuple) AmIDestination() (bool, error) {
	myHostIp, err := netutil.IsMyActiveHostIp(f.DestAddress)
	if err != nil {
		return false, err
	}
	return myHostIp, nil
}
