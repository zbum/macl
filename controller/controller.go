package controller

import (
	"encoding/json"
	"log/slog"
	controlsignal "macl/signal"
	"net"
	"strconv"
)

type Controller struct {
	controlPort int
	log         *slog.Logger
}

func NewController(log *slog.Logger, controlPort int) *Controller {
	return &Controller{
		controlPort: controlPort,
		log:         log,
	}
}

func (c *Controller) Start(fiveTuples []*controlsignal.FiveTuple) {

	for _, fiveTuple := range fiveTuples {
		var requestSignal = &controlsignal.ControlSignal{
			Command:   "1",
			FiveTuple: *fiveTuple,
		}

		c.sendRequestSignalToDestination(requestSignal)
		response := c.sendRequestSignalToSource(requestSignal)
		c.log.Info("[macl-controller]", "result", response)
	}
}

func (c *Controller) sendRequestSignalToDestination(requestSignal *controlsignal.ControlSignal) {
	var log = c.log
	destAddress := requestSignal.FiveTuple.DestAddress
	udpAddr, err := net.ResolveUDPAddr("udp", net.JoinHostPort(destAddress, strconv.Itoa(c.controlPort)))
	if err != nil {
		log.Error("[macl-controller] address failed", "dest_address", destAddress, "error", err)
		return
	}
	connection, err := net.DialUDP("udp", nil, udpAddr)
	defer connection.Close()
	if err != nil {
		log.Error("[macl-controller] connection failed", "error", err)
		return
	} else {
		log.Debug("[macl-controller] connection success")

		requestSignalBytes, err := json.Marshal(requestSignal)
		if err != nil {
			log.Warn("[macl-controller] unmarshal failed", "error", err)
			return
		}
		_, err = connection.Write(requestSignalBytes)
		if err != nil {
			log.Warn("[macl-controller] send failed")
			return
		} else {
			log.Info("[macl-controller] send success")
		}

		buffer := make([]byte, 1000)
		read, addr, err := connection.ReadFromUDP(buffer)
		if err != nil {
			return
		}

		log.Warn("[macl-controller] receive success", "localAddr", addr.String())
		log.Debug("[macl-controller] receive success", "payload", string(buffer[:read]))
	}
}

func (c *Controller) sendRequestSignalToSource(requestSignal *controlsignal.ControlSignal) *controlsignal.ResponseSignal {
	var log = c.log
	sourceAddress := requestSignal.FiveTuple.SrcAddress
	udpAddr, err := net.ResolveUDPAddr("udp", net.JoinHostPort(sourceAddress, strconv.Itoa(c.controlPort)))
	if err != nil {
		log.Error("[macl-controller] address failed", "dest_address", sourceAddress, "error", err)
		return controlsignal.NewFailResponseSignal(requestSignal.FiveTuple.TxId, "address failed", err)
	}
	connection, err := net.DialUDP("udp", nil, udpAddr)
	defer connection.Close()
	if err != nil {
		log.Error("[macl-controller] connection failed", "error", err)
		return controlsignal.NewFailResponseSignal(requestSignal.FiveTuple.TxId, "connection failed", err)
	} else {
		log.Debug("[macl-controller] connection success")

		requestSignalBytes, err := json.Marshal(requestSignal)
		if err != nil {
			log.Warn("[macl-controller] unmarshal failed", "error", err)
			return controlsignal.NewFailResponseSignal(requestSignal.FiveTuple.TxId, "unmarshal failed", err)
		}
		_, err = connection.Write(requestSignalBytes)
		if err != nil {
			log.Warn("[macl-controller] send failed")
			return controlsignal.NewFailResponseSignal(requestSignal.FiveTuple.TxId, "send failed", err)
		} else {
			log.Info("[macl-controller] send success")
		}

		buffer := make([]byte, 1000)
		read, addr, err := connection.ReadFromUDP(buffer)
		if err != nil {
			return controlsignal.NewFailResponseSignal(requestSignal.FiveTuple.TxId, "receive response failed", err)
		}

		log.Warn("[macl-controller] receive success", "localAddr", addr.String())
		log.Debug("[macl-controller] receive success", "payload", string(buffer[:read]))

		return controlsignal.NewSuccessResponseSignal(requestSignal.FiveTuple.TxId, &requestSignal.FiveTuple)
	}
}
