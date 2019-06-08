package data

import (
	"encoding/json"
)

type Data struct {
	Password      string  `json:"password,omitempty"`
	UUID          string  `json:"UUID"`
	IPv4Addr      string  `json:"IPv4Addr"`
	IPv6Addr      string  `json:"IPv6Addr"`
	ServerName    string  `json:"serverName"`
	System        string  `json:"system"`
	Location      string  `json:"location"`
	Uptime        string  `json:"uptime"`
	DownloadSpeed string  `json:"downloadSpeed"`
	UploadSpeed   string  `json:"uploadSpeed"`
	CPUUsed       int     `json:"CPUUsed"`
	MemoryUsed    float64 `json:"memoryUsed"`
	MemoryTotal   float64 `json:"memoryTotal"`
	SwapUsed      float64 `json:"swapUsed"`
	SwapTotal     float64 `json:"swapTotal"`
	DiskUsed      float64 `json:"diskUsed"`
	DiskTotal     float64 `json:"diskTotal"`
	LastTime      int64   `json:"lastTime"`
}

func EncryptData(d Data) ([]byte, error) {
	return json.Marshal(d)
}

func DecryptData(b []byte, d *Data) error {
	return json.Unmarshal(b, d)
}
