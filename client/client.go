package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/goomadao/serverstatus/client/send"
	"github.com/goomadao/serverstatus/client/status"
	"github.com/goomadao/serverstatus/util/data"
	"github.com/goomadao/serverstatus/util/logger"
)

var (
	server, serverName, password, location, logFile, logLevel *string
	port                                                      *int
	help                                                      *bool
)

func init() {
	server = flag.String("s", "", "The server address or domain")
	port = flag.Int("p", 36580, "The server port")
	serverName = flag.String("n", "", "Name of this server")
	password = flag.String("k", "", "Password to connect to server")
	location = flag.String("l", "", "Location of this server")
	logFile = flag.String("L", "./serverstatus.log", "The path to store the log file")
	logLevel = flag.String("level", "info", "Log levels: [ debug, info, error]")
	help = flag.Bool("h", false, "This help")

	flag.Usage = usage
}

func usage() {
	fmt.Fprintf(os.Stderr, `serverstatus client
Usage: serverstatus -s serverAddress [-p port(default 36580)] -L logFile -level logLevel -n serverName -k password -l location -h help

Options:
`)
	flag.PrintDefaults()
}

func main() {
	flag.Parse()
	status.Password = *password
	send.ServerAddr = *server
	send.ServerPort = strconv.Itoa(*port)
	status.ServerName = *serverName
	status.Location = *location
	logger.LogFile = *logFile
	logger.LogLevel = *logLevel

	if *help {
		flag.Usage()
		return
	}

	logger.InitLogger()

	var d *data.Data = new(data.Data)
	status.InitStatus()
	for {
		status.GetStatus(d)
		b, err := data.EncryptData(*d)
		if err != nil {
			fmt.Println(err)
			continue
		}
		go send.SendStatus(b)
		time.Sleep(2 * time.Second)
	}
}
