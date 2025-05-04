package agent

import (
	"encoding/json"
	"fmt"
	"log/slog"
	netutil "macl/net"
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

func NewAgent(controlPort int) *Agent {
	return &Agent{
		controlPort: controlPort,
		logger:      slog.New(slog.NewTextHandler(os.Stdout, nil)),
	}
}

func (a *Agent) Start() {
	startAgent(a.logger, a.controlPort)

	// Wait for a signal to terminate the server
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	<-done
}

func startAgent(log *slog.Logger, port int) {
	udpAddr, err := net.ResolveUDPAddr("udp", net.JoinHostPort("", strconv.Itoa(port)))
	if err != nil {
		log.Info("macl agent start failed", "error", err)
		return
	}
	connection, err := net.ListenUDP("udp4", udpAddr)
	if err != nil {
		log.Warn("macl agent start failed", "error", err)
		return
	}
	defer connection.Close()
	log.Info("macl agent started", "udpAddr", udpAddr)

	buffer := make([]byte, 8192)
	for {
		read, raddr, err := connection.ReadFromUDP(buffer)
		if err != nil {
			log.Warn("acl control signal 수신 실패 : %v\n", err)
		}
		signalBytes := buffer[:read]
		log.Info("acl control signal received successfully")
		log.Debug("acl control signal payload", string(signalBytes))

		controlSignal, err := parseControlSignal(signalBytes)
		if err != nil {
			log.Warn("acl control signal parse failed", err)
			sendResponse(connection, raddr, controlsignal.NewFailResponseSignal(0, "acl control signal parse 실패", err))
			continue
		}

		log.Info("acl control signal parse successfully")

		if isServer, err := amIServer(controlSignal.FiveTuple); err == nil && isServer {
			log.Info("acl control signal amIServer 성공 : %s\n", controlSignal)
			if controlSignal.FiveTuple.Protocol == "tcp" {
				NewTestReceiver(log).receivePacketFromSource(&controlSignal)
			}
		}
		if err != nil {
			log.Warn("acl control signal amIServer failed : %v\n", err)
			sendResponse(connection, raddr, controlsignal.NewFailResponseSignal(0, "acl control signal amIServer 실패", err))
			continue
		}

		if isClient, err := amIClient(controlSignal.FiveTuple); err == nil && isClient {
			log.Info("acl control signal amIClient 성공 : %s\n", controlSignal)
			if controlSignal.FiveTuple.Protocol == "tcp" {
				NewTestSender(log).sendPacketToDestination(&controlSignal)
			}
		}
		if err != nil {
			log.Warn("acl control signal amIClient 실패 : %v\n", err)
			sendResponse(connection, raddr, controlsignal.NewFailResponseSignal(0, "acl control signal amIClient 실패", err))
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

func amIClient(tuple controlsignal.FiveTuple) (bool, error) {
	myHostIp, err := netutil.IsMyActiveHostIp(tuple.SrcAddress)
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	return myHostIp, nil
}

func amIServer(tuple controlsignal.FiveTuple) (bool, error) {
	myHostIp, err := netutil.IsMyActiveHostIp(tuple.DestAddress)
	if err != nil {
		return false, err
	}
	return myHostIp, nil
}

func sendResponse(connection *net.UDPConn, addr *net.UDPAddr, responseSignal *controlsignal.ResponseSignal) {
	response, err := json.Marshal(responseSignal)
	if err != nil {
		fmt.Printf("응답 송신 실패 : %v\n", err)
	}

	_, err = connection.WriteToUDP(response, addr)
	if err != nil {
		fmt.Printf("응답 송신 실패 : %v\n", err)
	}
}
