package agent

import (
	"io"
	"log/slog"
	controlsignal "macl/signal"
	"net"
	"strconv"
	"time"
)

type TestReceiver struct {
	log *slog.Logger
}

func NewTestReceiver(log *slog.Logger) *TestReceiver {
	return &TestReceiver{
		log: log,
	}
}

func (r *TestReceiver) receivePacketFromSource(controlSignal *controlsignal.ControlSignal) {

	if controlSignal.FiveTuple.Protocol == "tcp" {
		go receiveTcpSignal(r.log, controlSignal.FiveTuple.DestPort)
	} else if controlSignal.FiveTuple.Protocol == "udp" {
		go receiveUdpSignal(r.log, controlSignal.FiveTuple.DestPort)
	} else {
		r.log.Warn("[macl-agent-receiver] unknown protocol", "protocol", controlSignal.FiveTuple.Protocol)
	}
}

func receiveUdpSignal(log *slog.Logger, port int) {
	udpAddr, err := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(port))
	if err != nil {
		log.Warn("[macl-agent-receiver] start failed", "error", err)
		return
	}

	connection, err := net.ListenUDP("udp4", udpAddr)
	if err != nil {
		log.Warn("[macl-agent-receiver] start failed", "error", err)
		return
	}
	defer connection.Close()
	log.Info("[macl-agent-receiver] started", "address", connection.RemoteAddr())

	buffer := make([]byte, 1000)

	err = connection.SetDeadline(time.Now().Add(5 * time.Second))
	if err != nil {
		log.Warn("[macl-agent-receiver] receive failed", "error", err)
		return
	}

	read, addr, err := connection.ReadFromUDP(buffer)
	if err != nil {
		if err == io.EOF {
			log.Warn("[macl-agent-receiver] receive failed", "error", err)
		} else {
			log.Warn("[macl-agent-receiver] receive failed", "error", err)
		}
		return
	}

	log.Info("[macl-agent-receiver] received successfully", "port", strconv.Itoa(port))
	log.Debug("[macl-agent-receiver] received payload", string(buffer[:read]))

	_, err = connection.WriteToUDP(buffer[:read], addr)
	if err != nil {
		log.Warn("[macl-agent-receiver] send failed : %v\n", err)
		return
	}

	log.Info("[macl-agent-receiver] send successfully")
}

func receiveTcpSignal(log *slog.Logger, port int) {

	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		log.Warn("[macl-agent-receiver] start failed", "error", err)
		return
	}
	defer listener.Close()
	log.Info("[macl-agent-receiver] started", "address", listener.Addr())

	connection, err := listener.Accept()
	if err != nil {
		return
	}
	defer connection.Close()

	log.Info("[macl-agent-receiver] connected ", "remote", connection.RemoteAddr(), "local", connection.LocalAddr())

	buffer := make([]byte, 1000)

	err = connection.SetDeadline(time.Now().Add(5 * time.Second))
	if err != nil {
		log.Warn("[macl-agent-receiver] receive failed", "error", err)
		return
	}
	read, err := connection.Read(buffer)
	if err != nil {
		if err == io.EOF {
			log.Warn("[macl-agent-receiver] receive failed", "error", err)
		} else {
			log.Warn("[macl-agent-receiver] receive failed", "error", err)
		}
		return
	}
	log.Info("[macl-agent-receiver] received successfully", "port", strconv.Itoa(port))
	log.Debug("[macl-agent-receiver] received payload", string(buffer[:read]))

	_, err = connection.Write(buffer[:read])
	if err != nil {
		log.Warn("[macl-agent-receiver] send failed : %v\n", err)
		return
	}

	log.Info("[macl-agent-receiver] send successfully")
}
