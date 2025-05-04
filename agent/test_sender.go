package agent

import (
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

func NewTestSender(log *slog.Logger) *TestSender {
	return &TestSender{
		log: log,
	}
}

func (s *TestSender) sendPacketToDestination(controlSignal *controlsignal.ControlSignal) {

	for trial := range 3 {
		s.log.Debug("trial", "count", trial)
		if controlSignal.FiveTuple.Protocol == "tcp" {
			_, err := sendTcpPacket(s.log, controlSignal.FiveTuple.DestJoinedAddress())
			if err != nil {
				s.log.Warn("sendTcpPacket", "error", err)
				return
			}
		}
		time.Sleep(1 * time.Second)
	}
}

func sendTcpPacket(log *slog.Logger, joinedAddress string) (bool, error) {
	connection, err := net.DialTimeout("tcp", joinedAddress, 1*time.Second)

	if err != nil {
		log.Warn("[macl-agent-sender] connection failed:", "error", err)
		return false, err
	} else {
		defer connection.Close()
		log.Info("[macl-agent-sender] connection success:", "joinedAddress", joinedAddress)

		localAddr := connection.LocalAddr().String()

		_, err := connection.Write([]byte(httpMessage))
		if err != nil {
			log.Info("송신 실패 : %v\n", err)
			return false, err
		} else {
			log.Info("[macl-agent-sender] send success:", "payload", httpMessage)
		}

		buffer := make([]byte, 50)
		read, err := connection.Read(buffer)
		if err != nil {
			log.Info("[macl-agent-sender] receive failed:", "localAddr", localAddr, "error", err)
			return false, err
		}
		log.Info("[macl-agent-sender] callback packet received successfully")
		log.Debug("[macl-agent-sender] callback packet", "payload", string(buffer[:read]))

		return true, nil
	}
}
