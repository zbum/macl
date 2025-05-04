package agent

import (
	"fmt"
	"io"
	controlsignal "macl/signal"
	"net"
	"strconv"
	"time"
)

func receivePacketFromSource(controlSignal *controlsignal.ControlSignal) {

	if controlSignal.FiveTuple.Protocol == "tcp" {
		go receiveTcpSignal(controlSignal.FiveTuple.DestPort)
	}
}

func receiveTcpSignal(port int) {

	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		fmt.Printf("macl receiver start failed: %v \n", err)
		return
	}
	fmt.Printf("macl receiver started: %v \n", listener.Addr())

	connection, err := listener.Accept()
	if err != nil {
		return
	}
	defer connection.Close()

	fmt.Printf("macl receiver connected : Remote: %v, Local: %v\n", connection.RemoteAddr(), connection.LocalAddr())

	buffer := make([]byte, 1000)

	err = connection.SetDeadline(time.Now().Add(5 * time.Second))
	if err != nil {
		fmt.Printf("macl server receive failed : %v\n", err)
		return
	}
	read, err := connection.Read(buffer)
	if err != nil {
		if err == io.EOF {
			fmt.Printf("macl server receive failed : %v\n", err)
		} else {
			fmt.Printf("macl server receive failed : %v\n", err)
		}
		return
	}

	fmt.Printf("macl server received successfully [%s]: %s\n", net.JoinHostPort("", strconv.Itoa(port)), string(buffer[:read]))
	_, err = connection.Write(buffer[:read])

}
