package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/goomadao/serverstatus/server/auth"
	"github.com/goomadao/serverstatus/server/dashboard"
	"github.com/goomadao/serverstatus/server/receive"
	"github.com/goomadao/serverstatus/server/status"
	"github.com/goomadao/serverstatus/util/logger"
	
	"go.uber.org/zap"
)

var (
	port, webPort                           *int
	statusFile, password, logFile, logLevel *string
	help                                    *bool
)

func init() {
	port = flag.Int("P", 36580, "The port to receive information from clients")
	webPort = flag.Int("p", 8080, "The port for the dashboard")
	statusFile = flag.String("f", "./status.json", "The position to store received status")
	password = flag.String("k", "", "Password to connect to this server")
	logFile = flag.String("L", "./serverstatus.log", "The path to store the log file")
	logLevel = flag.String("level", "info", "Log levels: [ debug, info, error]")
	help = flag.Bool("h", false, "This help")

	flag.Usage = usage
}

func usage() {
	fmt.Fprintf(os.Stderr, `serverstatus server
Usage: serverstatus -k password [-P udp port(default 36580)] [-p dashboard port(default 8080)] -L logFile -level logLevel [-f statusFile(default ./status.json] [-h help]
	
Options:
`)
	flag.PrintDefaults()
}

func main() {
	flag.Parse()
	receive.Port = *port
	dashboard.Port = *webPort
	status.StatusPath = *statusFile
	auth.Password = *password
	logger.LogFile = *logFile
	logger.LogLevel = *logLevel

	if *help {
		flag.Usage()
		return
	}
	logger.InitLogger()
	go dashboard.Dashboard()
	// logger.Logger.Info("Dashboard started!")
	err := receive.ListenPort()
	if err != nil {
		logger.Logger.Fatal("Listen to [::]:"+strconv.Itoa(*port)+" failed",
			zap.Error(err))
	}
}
