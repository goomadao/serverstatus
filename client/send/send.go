package send

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"os"

	"github.com/goomadao/serverstatus/util/logger"

	"go.uber.org/zap"
)

var (
	ServerAddr string
	ServerPort string
)

func SendStatus(d []byte) {
	socket := ServerAddr + ":" + ServerPort
	udpAddr, err := net.ResolveUDPAddr("udp", socket)
	if err != nil {
		logger.Logger.Error("resolve udp addr failed",
			zap.Error(err))
		return
	}
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		logger.Logger.Error("dial udp addr failed",
			zap.Error(err))
		return
	}
	_, err = conn.Write(d)
	if err != nil {
		logger.Logger.Error("send udp failed",
			zap.Error(err))
		return
	}
	// printJson(d)
	conn.Close()
	logger.Logger.Info("Send data success!")
}

func printJson(d []byte) {
	var out bytes.Buffer
	err := json.Indent(&out, d, "", "\t")

	if err != nil {
		fmt.Println("format json failed")
		return
	}
	out.WriteTo(os.Stdout)
}
