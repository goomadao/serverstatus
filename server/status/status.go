package status

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"
	"time"

	"github.com/goomadao/serverstatus/server/auth"
	"github.com/goomadao/serverstatus/util/data"
	"github.com/goomadao/serverstatus/util/logger"
	
	"go.uber.org/zap"
)

var (
	Mu         sync.Mutex
	StatusPath string
)

type Servers struct {
	Servers    []data.Data `json:"servers"`
	TmpServers []data.Data `json:"tmpServers"`
}

func HandleData(d []byte) {
	Mu.Lock()
	defer Mu.Unlock()
	var newData *data.Data = new(data.Data)
	err := data.DecryptData(d, newData)
	if err != nil {
		logger.Logger.Warn("Wrong json data",
			zap.Error(err))
		return
	}
	if newData.Password != auth.Password {
		logger.Logger.Warn("Wrong password!")
		return
	}
	var allServers *Servers = new(Servers)
	getServers(allServers)
	flag := true
	for index := 0; index < len(allServers.Servers); index++ {
		status := allServers.Servers[index]
		if flag && status.UUID == newData.UUID {
			updateServer(&allServers.Servers[index], newData)
			allServers.Servers[index].LastTime = time.Now().Unix()
			flag = false
		}
	}
	for index := 0; index < len(allServers.TmpServers); index++ {
		status := allServers.TmpServers[index]
		if flag && status.UUID == newData.UUID {
			updateServer(&allServers.TmpServers[index], newData)
			allServers.TmpServers[index].LastTime = time.Now().Unix()
			flag = false
		} else if time.Now().Unix()-allServers.TmpServers[index].LastTime >= 60*60*24*3 {
			allServers.TmpServers = append(allServers.TmpServers[:index], allServers.TmpServers[index+1:]...)
			index--
		}
	}
	if flag {
		newData.LastTime = time.Now().Unix()
		newData.Password = ""
		allServers.TmpServers = append(allServers.TmpServers, *newData)
	}
	saveServers(allServers)
}

func getServers(servers *Servers) {
	_, err := os.Stat(StatusPath)
	if err != nil {
		if os.IsNotExist(err) {
			_, err2 := os.Create(StatusPath)
			if err2 != nil {
				logger.Logger.Error("Create status file failed",
					zap.Error(err2))
				return
			}
		} else {
			logger.Logger.Error("Get stat of status file failed",
				zap.Error(err))
			return
		}
	}
	file, err := os.Open(StatusPath)
	if err != nil {
		logger.Logger.Error("Open status file failed",
			zap.Error(err))
		return
	}
	defer file.Close()
	allServers, err := ioutil.ReadAll(file)
	if err != nil {
		logger.Logger.Error("read status file failed",
			zap.Error(err))
		return
	}
	err = json.Unmarshal(allServers, servers)
	if err != nil {
		logger.Logger.Error("Turn json file into struct failed",
			zap.Error(err))
		return
	}
}

func GetServers(servers *Servers) {
	Mu.Lock()
	defer Mu.Unlock()
	_, err := os.Stat(StatusPath)
	if err != nil {
		if os.IsNotExist(err) {
			_, err2 := os.Create(StatusPath)
			if err2 != nil {
				logger.Logger.Error("Create status file failed",
					zap.Error(err2))
				return
			}
		} else {
			logger.Logger.Error("Get stat of status file failed",
				zap.Error(err))
			return
		}
	}
	file, err := os.Open(StatusPath)
	if err != nil {
		logger.Logger.Error("Open status file failed",
			zap.Error(err))
		return
	}
	defer file.Close()
	allServers, err := ioutil.ReadAll(file)
	if err != nil {
		logger.Logger.Error("read status file failed",
			zap.Error(err))
		return
	}
	err = json.Unmarshal(allServers, servers)
	if err != nil {
		logger.Logger.Error("Turn json file into struct failed",
			zap.Error(err))
		return
	}
}

func updateServer(server, newData *data.Data) {
	server.IPv4Addr = newData.IPv4Addr
	server.IPv6Addr = newData.IPv6Addr
	server.ServerName = newData.ServerName
	server.System = newData.System
	server.Location = newData.Location
	server.Uptime = newData.Uptime
	server.DownloadSpeed = newData.DownloadSpeed
	server.UploadSpeed = newData.UploadSpeed
	server.CPUUsed = newData.CPUUsed
	server.MemoryUsed = newData.MemoryUsed
	server.MemoryTotal = newData.MemoryTotal
	server.SwapUsed = newData.SwapUsed
	server.SwapTotal = newData.SwapTotal
	server.DiskUsed = newData.DiskUsed
	server.DiskTotal = newData.DiskTotal
}

func saveServers(servers *Servers) {
	file, err := os.OpenFile(
		StatusPath,
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		0666,
	)
	if err != nil {
		logger.Logger.Error("open status file failed",
			zap.Error(err))
		return
	}
	defer file.Close()
	allServers, err := json.Marshal(*servers)
	if err != nil {
		logger.Logger.Error("struct to []byte failed",
			zap.Error(err))
		return
	}
	var out bytes.Buffer
	err = json.Indent(&out, allServers, "", "\t")
	if err != nil {
		logger.Logger.Error("json.Indent failed",
			zap.Error(err))
		return
	}
	_, err = out.WriteTo(file)
	if err != nil {
		logger.Logger.Error("Write to status file failed",
			zap.Error(err))
		return
	}
}

func MoveServer(from, UUID, action string) bool {
	Mu.Lock()
	defer Mu.Unlock()
	allServers := new(Servers)
	getServers(allServers)
	if from == "temp" {
		if action == "save" {
			for index, value := range allServers.TmpServers {
				if value.UUID == UUID {
					allServers.Servers = append(allServers.Servers, value)
					allServers.TmpServers = append(allServers.TmpServers[:index], allServers.TmpServers[index+1:]...)
					saveServers(allServers)
					logger.Logger.Info("Save success!")
					return true
				}
			}
		} else if action == "delete" {
			for index, value := range allServers.TmpServers {
				if value.UUID == UUID {
					allServers.TmpServers = append(allServers.TmpServers[:index], allServers.TmpServers[index+1:]...)
					saveServers(allServers)
					logger.Logger.Info("Delete success!")
					return true
				}
			}
		}
	} else if from == "saved" {
		for index, value := range allServers.Servers {
			if value.UUID == UUID {
				allServers.TmpServers = append(allServers.TmpServers, value)
				allServers.Servers = append(allServers.Servers[:index], allServers.Servers[index+1:]...)
				saveServers(allServers)
				logger.Logger.Info("Remove success!")
				return true
			}
		}
	}
	return false
}
