package status

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/goomadao/serverstatus/util/data"
	"github.com/goomadao/serverstatus/util/logger"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"

	"go.uber.org/zap"
)

type traffic struct {
	receive, send uint64
	t             time.Time
}

var (
	lastTraffic       traffic
	Password          string
	UUID              string
	ServerName        string
	System            string
	Location          string
	trafficChan       chan traffic = make(chan traffic)
	UUIDChan          chan string  = make(chan string)
	IPv4AddrChan      chan string  = make(chan string)
	IPv6AddrChan      chan string  = make(chan string)
	ServerNameChan    chan string  = make(chan string)
	SystemChan        chan string  = make(chan string)
	LocationChan      chan string  = make(chan string)
	UptimeChan        chan string  = make(chan string)
	DownloadSpeedChan chan string  = make(chan string)
	UploadSpeedChan   chan string  = make(chan string)
	CPUUsageChan      chan int     = make(chan int)
	MemoryUsedChan    chan float64 = make(chan float64)
	MemoryTotalChan   chan float64 = make(chan float64)
	SwapUsedChan      chan float64 = make(chan float64)
	SwapTotalChan     chan float64 = make(chan float64)
	DiskUsedChan      chan float64 = make(chan float64)
	DiskTotalChan     chan float64 = make(chan float64)
)

func InitStatus() {
	go getUUID(UUIDChan)
	go getSystem(SystemChan)
	go getTraffic(trafficChan)
	UUID = <-UUIDChan
	System = <-SystemChan
	lastTraffic = <-trafficChan
	logger.Logger.Info("Init status success!")
}

func GetStatus(d *data.Data) {
	go getIPv4Addr(IPv4AddrChan)
	go getIPv6Addr(IPv6AddrChan)
	go getUptime(UptimeChan)
	go getNetSpeed(trafficChan, DownloadSpeedChan, UploadSpeedChan)
	go getCPUUsed(CPUUsageChan)
	go getMemory(MemoryTotalChan, MemoryUsedChan, SwapTotalChan, SwapUsedChan)
	go getDisk(DiskTotalChan, DiskUsedChan)
	d.IPv4Addr = <-IPv4AddrChan
	d.IPv6Addr = <-IPv6AddrChan
	d.Uptime = <-UptimeChan
	d.DownloadSpeed = <-DownloadSpeedChan
	d.UploadSpeed = <-UploadSpeedChan
	d.CPUUsed = <-CPUUsageChan
	d.MemoryTotal = <-MemoryTotalChan
	d.MemoryUsed = <-MemoryUsedChan
	d.SwapTotal = <-SwapTotalChan
	d.SwapUsed = <-SwapUsedChan
	d.DiskTotal = <-DiskTotalChan
	d.DiskUsed = <-DiskUsedChan
	d.Password, d.UUID, d.ServerName, d.System, d.Location = Password, UUID, ServerName, System, Location
	logger.Logger.Info("Get status success!")
}

func getTraffic(c chan traffic) {
	if runtime.GOOS == "darwin" {
		c <- *new(traffic)
		logger.Logger.Warn("Getting traffic in darwin isn't implemented till now. Return 0.")
		return
	} else if runtime.GOOS == "linux" {
		cmd := exec.Command("bash", "-c", ` cat /proc/net/dev | awk 'BEGIN{receive=0;send=0}{if ($1 !~ /Inter-/ && $1 !~ /face/ && $1 !~ /vir/ && $1 !~ /lo/ && $1 !~ /docker/){receive+=$2;send+=$10}}END{print receive,send}'`)
		time := time.Now()
		b, err := cmd.Output()
		if err != nil {
			logger.Logger.Warn("Get traffic from /proc/net/dev failed. Return nil struct traffic.",
				zap.Error(err))
			c <- *new(traffic)
			return
		}
		receiveAndSend := strings.Split(string(b[:len(b)-1]), " ")
		receive, err := strconv.ParseUint(receiveAndSend[0], 10, 64)
		if err != nil {
			logger.Logger.Warn("Get receive from string failed. Return nil struct traffic.",
				zap.Error(err))
			c <- *new(traffic)
			return
		}
		send, err := strconv.ParseUint(receiveAndSend[1], 10, 64)
		if err != nil {
			logger.Logger.Warn("Get send from string failed. Return nil struct traffic.",
				zap.Error(err))
			c <- *new(traffic)
			return
		}
		var ret traffic
		ret.receive, ret.send, ret.t = receive, send, time
		c <- ret
	}
}

func getUUID(c chan string) {
	os := runtime.GOOS
	if os == "darwin" {
		cmd := exec.Command("bash", "-c", `system_profiler SPHardwareDataType | awk '/UUID/ {print $3}'`)
		for {
			b, err := cmd.Output()
			if err != nil {
				logger.Logger.Error("Get product UUID failed!",
					zap.Error(err))
				continue
			}
			if b[len(b)-1] == 10 {
				c <- string(b[:len(b)-1])
			} else {
				c <- string(b)
			}
			break
		}
	} else if os == "linux" {
		for {
			uuid, err := ioutil.ReadFile("/sys/class/dmi/id/product_uuid")
			if err != nil {
				logger.Logger.Error("Get product_uuid failed",
					zap.Error(err))
				continue
			}
			if uuid[len(uuid)-1] == 10 {
				c <- string(uuid[:len(uuid)-1])
			} else {
				c <- string(uuid)
			}
		}

	}
}

func getIPv4Addr(c chan string) {
	cli := &http.Client{
		Timeout: 2 * time.Second,
	}
	r, err := cli.Get("http://v4.ipv6-test.com/api/myip.php")
	if err != nil {
		logger.Logger.Warn("Fetch ipv4 addr from http://v4.ipv6-test.com/api/myip.php failed",
			zap.Error(err))
		c <- ""
		return
	}
	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Logger.Warn("Read from http response Body failed",
			zap.Error(err))
		c <- ""
		return
	}
	c <- string(b)
}

func getIPv6Addr(c chan string) {
	cli := &http.Client{
		Timeout: 2 * time.Second,
	}
	r, err := cli.Get("http://v6.ipv6-test.com/api/myip.php")
	if err != nil {
		logger.Logger.Warn("Fetch ipv6 addr from http://v6.ipv6-test.com/api/myip.php failed",
			zap.Error(err))
		c <- ""
		return
	}
	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Logger.Warn("Read from http response Body failed",
			zap.Error(err))
		c <- ""
		return
	}
	c <- string(b)
}

func getSystem(c chan string) {
	if runtime.GOOS == "darwin" {
		c <- runtime.GOOS + "\n" + runtime.GOARCH
		return
	}
	file, err := os.Open("/etc/os-release")
	if err != nil {
		logger.Logger.Warn("/etc/os-release not exists!")
		c <- runtime.GOOS + " " + runtime.GOARCH
		return
	}
	defer file.Close()
	buff := bufio.NewReader(file)
	line, err := buff.ReadString('\n')
	if err != nil {
		logger.Logger.Error("Read from /etc/os-release failed. Return linux instead.",
			zap.Error(err))
		c <- runtime.GOOS + " " + runtime.GOARCH
		return
	}
	line = line[strings.Index(line, "=")+1 : len(line)-1]
	line = strings.ReplaceAll(line, "\"", "")
	line = strings.Trim(line, " ")
	if strings.Index(line, " ") != -1 {
		line = line[:strings.Index(line, " ")]
	}
	line = strings.Trim(line, " ")
	c <- line + "\n" + runtime.GOARCH
}

func getUptime(c chan string) {
	// Get uptime from file in linux
	// file, err := os.Open("/proc/uptime")
	// if err != nil {
	// 	fmt.Println(err)
	// 	panic("Open /proc/uptime failed")
	// 	return
	// }
	// defer file.Close()
	// buff := bufio.NewReader(file)
	// uptime, err := buff.ReadString(' ')
	// if err != nil {
	// 	fmt.Println(err)
	// 	panic("Read /proc/uptime failed")
	// 	return
	// }
	// uptime = strings.ReplaceAll(uptime, " ", "s")
	// duration, err := time.ParseDuration(uptime)
	// if err != nil {
	// 	fmt.Println(err)
	// 	panic("Parse to duration failed")
	// 	return
	// }
	uptime, err := host.Uptime()
	if err != nil {
		logger.Logger.Warn("Get uptime from package gopsutil failed",
			zap.Error(err))
		c <- "-"
		return
	}
	duration, err := time.ParseDuration(strconv.FormatUint(uptime, 10) + "s")
	if err != nil {
		logger.Logger.Warn("Parse to duration failed",
			zap.Error(err))
		c <- "-"
		return
	}
	day := int(duration.Hours()) / 24
	hour := int(duration.Hours()) % 24
	minute := int(duration.Minutes()) % 60
	second := int(duration.Seconds()) % 60
	s := ""
	if day > 1 {
		s = fmt.Sprintf("%ddays\n%02d:%02d:%02d", day, hour, minute, second)
	} else {
		s = fmt.Sprintf("%dday\n%02d:%02d:%02d", day, hour, minute, second)
	}
	c <- s
}

func getNetSpeed(trafficChan chan traffic, DownloadSpeedChan, UploadSpeedChan chan string) {
	if runtime.GOOS == "darwin" {
		DownloadSpeedChan <- "-"
		UploadSpeedChan <- "-"
		return
	} else if runtime.GOOS == "linux" {
		go getTraffic(trafficChan)
		currentTraffic := <-trafficChan
		if lastTraffic.t.IsZero() {
			lastTraffic = currentTraffic
			DownloadSpeedChan <- "-"
			UploadSpeedChan <- "-"
			return
		} else {
			download := float64(currentTraffic.receive-lastTraffic.receive) / (currentTraffic.t.Sub(lastTraffic.t).Seconds())
			upload := float64(currentTraffic.send-lastTraffic.send) / (currentTraffic.t.Sub(lastTraffic.t).Seconds())
			lastTraffic = currentTraffic
			DownloadSpeedChan <- fmt.Sprintf("%.2f MB/s", download/1024/1024+0.005)
			UploadSpeedChan <- fmt.Sprintf("%.2f MB/s", upload/1024/1024+0.005)
			return
		}
	}
}

func getCPUUsed(c chan int) {
	// if runtime.GOOS == "darwin" {
	// 	pscmd := fmt.Sprintf(`ps -A -o %%cpu | awk '{s+=$1} END {print s/%d "%%"}'`, runtime.NumCPU())
	// 	cpu := exec.Command("bash", "-c", pscmd)
	// 	b, err := cpu.Output()
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		fmt.Println("Get cpu usage from ps failed")
	// 		c <- 0
	// 		return
	// 	}
	// 	usage, err := strconv.ParseFloat(string(b[:len(b)-1]), 64)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		fmt.Println("Parse cpu usage from string to float64 failed")
	// 		c <- 0
	// 		return
	// 	}
	// 	c <- int(usage)
	// } else if runtime.GOOS == "linux" {
	// 	u := cpuusage.Usage{}
	// 	err := u.Measure()
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		panic("get cpu usage failed")
	// 		c <- 0
	// 		return
	// 	}
	// 	c <- u.Overall
	// }
	cpuusage, err := cpu.Percent(time.Second, false)
	if err != nil {
		logger.Logger.Warn("Get cpu usage from package gopsutil failed. Return 0.",
			zap.Error(err))
		c <- 0
		return
	}
	c <- int(cpuusage[0])
}

func getMemory(MemoryTotalChan, MemoryUsedChan, SwapTotalChan, SwapUsedChan chan float64) {
	m, err := mem.VirtualMemory()
	if err != nil {
		logger.Logger.Warn("Get memory information from package gopsutil failed. Return 0.",
			zap.Error(err))
		MemoryTotalChan <- 0
		MemoryUsedChan <- 0
	} else {
		MemoryTotalChan <- float64(m.Total) / 1024 / 1024 / 1024
		MemoryUsedChan <- float64(m.Used) / 1024 / 1024 / 1024
	}
	s, err := mem.SwapMemory()
	if err != nil {
		logger.Logger.Warn("Get swap information from package gopsutil failed",
			zap.Error(err))
		SwapTotalChan <- 0
		SwapUsedChan <- 0
	} else {
		SwapTotalChan <- float64(s.Total) / 1024 / 1024 / 1024
		SwapUsedChan <- float64(s.Used) / 1024 / 1024 / 1024
	}
}

// Get memory and swap information from /proc/meminfo file in linux
// func getMemory(d *data.Data) {
// 	var total, used, stotal, sused float64
// 	meminfo, err := ioutil.ReadFile("/proc/meminfo")
// 	if err != nil {
// 		fmt.Println(err)
// 		panic("read meminfo failed")
// 		return
// 	}
// 	info := string(meminfo)
// 	info = strings.ReplaceAll(info, " ", "")
// 	info = strings.ReplaceAll(info, "\n", "")
// 	total, err = strconv.ParseFloat(info[strings.Index(info, "MemTotal")+9:strings.Index(info, "kBMemFree")], 64)
// 	if err != nil {
// 		fmt.Println(err)
// 		panic("get total memory failed")
// 		return
// 	}
// 	free, err := strconv.ParseFloat(info[strings.Index(info, "MemFree")+8:strings.Index(info, "kBMemAvailable")], 64)
// 	if err != nil {
// 		fmt.Println(err)
// 		panic("get free memory failed")
// 		return
// 	}
// 	buffers, err := strconv.ParseFloat(info[strings.Index(info, "Buffers")+8:strings.Index(info, "kBCached")], 64)
// 	if err != nil {
// 		fmt.Println(err)
// 		panic("get buffers memory failed")
// 		return
// 	}
// 	cached, err := strconv.ParseFloat(info[strings.Index(info, "Cached")+7:strings.Index(info, "kBSwapCached")], 64)
// 	if err != nil {
// 		fmt.Println(err)
// 		panic("get cached memory failed")
// 		return
// 	}
// 	sreclaimable, err := strconv.ParseFloat(info[strings.Index(info, "SReclaimable")+13:strings.Index(info, "kBSUnreclaim")], 64)
// 	if err != nil {
// 		fmt.Println(err)
// 		panic("get strclaimable memory failed")
// 		return
// 	}
// 	used = total - free - buffers - cached - sreclaimable
// 	stotal, err = strconv.ParseFloat(info[strings.Index(info, "SwapTotal")+10:strings.Index(info, "kBSwapFree")], 64)
// 	if err != nil {
// 		fmt.Println(err)
// 		panic("get swap total memory failed")
// 		return
// 	}
// 	sfree, err := strconv.ParseFloat(info[strings.Index(info, "SwapFree")+9:strings.Index(info, "kBDirty")], 64)
// 	if err != nil {
// 		fmt.Println(err)
// 		panic("get swap free memory failed")
// 		return
// 	}
// 	sused = stotal - sfree
// 	d.MemoryTotal = float64(total) / 1024 / 1024
// 	d.MemoryUsed = float64(used) / 1024 / 1024
// 	d.SwapTotal = float64(stotal) / 1024 / 1024
// 	d.SwapUsed = float64(sused) / 1024 / 1024
// }

func getDisk(total, used chan float64) {
	fs := syscall.Statfs_t{}
	err := syscall.Statfs("/", &fs)
	if err != nil {
		logger.Logger.Warn("disk syscall failed. Return 0.",
			zap.Error(err))
		total <- 0.0
		used <- 0.0
		return
	}
	all := fs.Blocks * uint64(fs.Bsize)
	free := fs.Bfree * uint64(fs.Bsize)
	total <- float64(all) / 1024 / 1024 / 1024
	used <- float64(all-free) / 1024 / 1024 / 1024

}
