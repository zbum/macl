package agent

import (
	"fmt"
	"log/slog"
	controlsignal "macl/signal"
	"net"
	"time"
)

const (
	httpMessage = `GET / HTTP/1.1
Host: localhost
User-Agent: aclchecker
Accept: */*

`
)

type TestSender struct {
	log *slog.Logger
}

func NewTestSender() *TestSender {
	return &TestSender{
		log: slog.Default(),
	}
}

func (s *TestSender) sendPacketToDestination(controlSignal *controlsignal.ControlSignal) {

	for trial := range 3 {
		fmt.Printf("trial %d\n", trial)
		if controlSignal.FiveTuple.Protocol == "tcp" {
			sendTcpPacket(s.log, controlSignal.FiveTuple.DestJoinedAddress())
		}
		time.Sleep(1 * time.Second)
	}
}

func sendTcpPacket(log *slog.Logger, joinedAddress string) (bool, error) {
	connection, err := net.DialTimeout("tcp", joinedAddress, 1*time.Second)

	if err != nil {
		log.Warn("연결 실패 : %v\n", err)
		return false, err
	} else {
		defer connection.Close()
		log.Info("연결 성공 : %s \n", joinedAddress)

		localAddr := connection.LocalAddr().String()

		_, err := connection.Write([]byte(httpMessage))
		if err != nil {
			log.Info("송신 실패 : %v\n", err)
			return false, err
		} else {
			log.Info("송신 성공 : %v\n", httpMessage)
		}

		buffer := make([]byte, 50)
		read, err := connection.Read(buffer)
		if err != nil {
			log.Info("수신 실패 [%s] %v\n", localAddr, err)
			return false, err
		}
		log.Info("macl packet received successfully : \n")
		log.Debug("macl packet : %s\n", string(buffer[:read]))

		return true, nil
	}
}
