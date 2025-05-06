package agent

import (
	"encoding/json"
	"log/slog"
	controlsignal "macl/signal"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

type Agent struct {
	controlPort int
	logger      *slog.Logger
}

func NewAgent(log *slog.Logger, controlPort int) *Agent {
	return &Agent{
		controlPort: controlPort,
		logger:      log,
	}
}

func (a *Agent) Start() {
	a.startUdpListener()

	// Wait for a signal to terminate the server
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	<-done
}

func (a *Agent) startUdpListener() {
	log := a.logger
	port := a.controlPort

	udpAddr, err := net.ResolveUDPAddr("udp", net.JoinHostPort("", strconv.Itoa(port)))
	if err != nil {
		log.Info("[macl-agent] start failed", "error", err)
		return
	}
	connection, err := net.ListenUDP("udp4", udpAddr)
	if err != nil {
		log.Warn("[macl-agent]start failed", "error", err)
		return
	}
	defer connection.Close()
	log.Info("[macl-agent] agent started", "udpAddr", udpAddr)

	a.requestProcess(connection)

}

func (a *Agent) requestProcess(connection *net.UDPConn) {
	log := a.logger

	buffer := make([]byte, 8192)
	for {
		read, raddr, err := connection.ReadFromUDP(buffer)
		if err != nil {
			log.Warn("[macl-agent] control signal read failed", "error", err)
		}
		signalBytes := buffer[:read]
		log.Info("[macl-agent] control signal received successfully")
		log.Debug("[macl-agent] control signal", "payload", string(signalBytes))

		controlSignal, err := parseControlSignal(signalBytes)
		if err != nil {
			log.Warn("[macl-agent] control signal parse failed", err)
			a.sendResponse(connection, raddr, controlsignal.NewFailResponseSignal("", "acl control signal parse 실패", err))
			continue
		}

		log.Info("[macl-agent] acl control signal parse successfully")

		if isDestination, err := controlSignal.FiveTuple.AmIDestination(); err == nil && isDestination {
			log.Info("[macl-agent] control signal amIServer 성공 : %s\n", controlSignal)
			NewTestReceiver(log).receivePacketFromSource(&controlSignal)
		}
		if err != nil {
			log.Warn("[macl-agent] control signal amIServer failed : %v\n", err)
			a.sendResponse(connection, raddr, controlsignal.NewFailResponseSignal(controlSignal.FiveTuple.TxId, "acl control signal amIServer 실패", err))
			continue
		}

		if isSource, err := controlSignal.FiveTuple.AmISource(); err == nil && isSource {
			log.Info("[macl-agent] control signal AmISource success", "signal", controlSignal)
			err := NewTestSender(log).sendPacketToDestination(&controlSignal)
			if err != nil {
				a.sendResponse(connection, raddr, controlsignal.NewFailResponseSignal(controlSignal.FiveTuple.TxId, "[macl-agent-sender] test failed", err))
				continue
			} else {
				a.sendResponse(connection, raddr, controlsignal.NewSuccessResponseSignal(controlSignal.FiveTuple.TxId, &controlSignal.FiveTuple))
			}
		}
		if err != nil {
			log.Warn("[macl-agent] control signal AmISource failed", "error", err)
			a.sendResponse(connection, raddr, controlsignal.NewFailResponseSignal(controlSignal.FiveTuple.TxId, "acl control signal amIClient 실패", err))
			continue
		}

	}
}

func parseControlSignal(signalBytes []byte) (controlsignal.ControlSignal, error) {
	var controlSignal controlsignal.ControlSignal
	err := json.Unmarshal(signalBytes, &controlSignal)
	if err != nil {
		return controlsignal.ControlSignal{}, err
	}
	return controlSignal, nil
}

func (a *Agent) sendResponse(connection *net.UDPConn, addr *net.UDPAddr, responseSignal *controlsignal.ResponseSignal) {
	log := a.logger

	response, err := json.Marshal(responseSignal)
	if err != nil {
		log.Warn("[macl-agent] failed to send response", "error", err)
	}

	_, err = connection.WriteToUDP(response, addr)
	if err != nil {
		log.Warn("[macl-agent] failed to send response", "error", err)
	}
}
