package receive

import (
	"net"
	"strconv"

	"github.com/goomadao/serverstatus/server/status"
	"github.com/goomadao/serverstatus/util/logger"
	
	"go.uber.org/zap"
)

var (
	Port int
)

func ListenPort() error {
	udpAddr, err := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(Port))
	if err != nil {
		logger.Logger.Error("Resolve udp addr failed",
			zap.Error(err))
		return err
	}
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		logger.Logger.Error(("Listen udp add failed"),
			zap.Error(err))
		return err
	}
	for {
		var buf [512]byte
		n, _, err := conn.ReadFromUDP(buf[0:])
		if err != nil {
			logger.Logger.Warn("Read from udp conn failed",
				zap.Error(err))
			continue
		}
		go status.HandleData(buf[0:n])
	}
}
