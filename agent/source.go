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

func NewTestSender(log *slog.Logger) *TestSender {
	return &TestSender{
		log: log,
	}
}

func (s *TestSender) sendPacketToDestination(controlSignal *controlsignal.ControlSignal) error {
	var err error
	for trial := range 3 {
		s.log.Debug("trial", "count", trial)
		if controlSignal.FiveTuple.Protocol == "tcp" {
			_, err = sendTcpPacket(s.log, controlSignal.FiveTuple.DestJoinedAddress())
			if err != nil {
				s.log.Warn("[macl-agent-sender] sendTcpPacket", "error", err)
				continue
			} else {
				break
			}
		} else if controlSignal.FiveTuple.Protocol == "udp" {
			_, err = sendUdpPacket(s.log, controlSignal.FiveTuple.DestJoinedAddress())
			if err != nil {
				s.log.Warn("[macl-agent-sender] sendUdpPacket", "error", err)
				continue
			} else {
				break
			}
		} else {
			s.log.Warn("[macl-agent-sender] unknown protocol", "protocol", controlSignal.FiveTuple.Protocol)
			return nil
		}
		time.Sleep(1 * time.Second)
	}
	return err
}

func sendUdpPacket(log *slog.Logger, joinedAddress string) (bool, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", joinedAddress)
	if err != nil {
		log.Warn("[macl-agent-sender] address failed:", "error", err)
		return false, err
	}
	connection, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		log.Warn("[macl-agent-sender] connection failed:", "error", err)
		return false, err
	} else {
		defer connection.Close()
		log.Info("[macl-agent-sender] connection success:", "joinedAddress", joinedAddress)

		localAddr := connection.LocalAddr().String()

		_, err := connection.Write([]byte(httpMessage))
		if err != nil {
			fmt.Printf("송신 실패 : %v\n", err)
			return false, err
		} else {
			log.Info("[macl-agent-sender] send success", "address", joinedAddress)
			log.Debug("[macl-agent-sender] send success", "payload", httpMessage)
		}

		err = connection.SetReadDeadline(time.Now().Add(3 * time.Second))
		if err != nil {
			log.Info("[macl-agent-sender] set deadline failed", "localAddr", localAddr, "error", err)
			return false, err
		}

		buffer := make([]byte, 1000)
		read, addr, err := connection.ReadFromUDP(buffer)
		if err != nil {
			log.Info("[macl-agent-sender] receive failed:", "remoteAddr", addr, "localAddr", localAddr, "error", err)
			return false, err
		}
		log.Info("[macl-agent-sender] callback packet received successfully")
		log.Debug("[macl-agent-sender] callback packet", "payload", string(buffer[:read]))

		return true, nil
	}

}

func sendTcpPacket(log *slog.Logger, joinedAddress string) (bool, error) {
	connection, err := net.DialTimeout("tcp", joinedAddress, 1*time.Second)

	if err != nil {
		log.Warn("[macl-agent-sender] connection failed", "error", err)
		return false, err
	} else {
		defer connection.Close()
		log.Info("[macl-agent-sender] connection success", "joinedAddress", joinedAddress)

		localAddr := connection.LocalAddr().String()

		_, err := connection.Write([]byte(httpMessage))
		if err != nil {
			log.Info("[macl-agent-sender] send failed", "error", err)
			return false, err
		} else {
			log.Info("[macl-agent-sender] send success", "address", joinedAddress)
			log.Debug("[macl-agent-sender] send success", "payload", httpMessage)
		}

		err = connection.SetReadDeadline(time.Now().Add(3 * time.Second))
		if err != nil {
			log.Info("[macl-agent-sender] set deadline failed", "localAddr", localAddr, "error", err)
			return false, err
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
