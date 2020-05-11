package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"io"
	"math"
	"os"
	"strings"
	"time"
)
type TimesStat struct {
	CPU       string  `json:"cpu"`
	User      float64 `json:"user"`
	System    float64 `json:"system"`
	Idle      float64 `json:"idle"`
	Nice      float64 `json:"nice"`
	Iowait    float64 `json:"iowait"`
	Irq       float64 `json:"irq"`
	Softirq   float64 `json:"softirq"`
	Steal     float64 `json:"steal"`
	Guest     float64 `json:"guest"`
	GuestNice float64 `json:"guestNice"`
}
//读取key=value类型的配置文件
func InitConfig(path string) map[string]string {
	config := make(map[string]string)

	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		panic(err)
	}
	r := bufio.NewReader(f)
	for {
		b, _, err := r.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		s := strings.TrimSpace(string(b))
		index := strings.Index(s, "=")
		if index < 0 {
			continue
		}
		key := strings.TrimSpace(s[:index])
		if len(key) == 0 {
			continue
		}
		value := strings.TrimSpace(s[index+1:])
		if len(value) == 0 {
			continue
		}
		config[key] = value
	}
	return config
}
func getCpuTime() ([]cpu.TimesStat, error) {
	res, err := cpu.Times(true) // false是展示全部总和 true是分布展示
	if err != nil {
		return []cpu.TimesStat{}, err
	}
	return res, nil
}
func getTotalCpuUsage() ([]float64, error) {
    return cpu.Percent(time.Second * 1, false)
}
func getAllBusy(t cpu.TimesStat) (float64, float64) {
	busy := t.User + t.System + t.Nice + t.Iowait + t.Irq +
		t.Softirq + t.Steal
	return busy + t.Idle, busy
}

func calculateBusy(t1, t2 cpu.TimesStat) float64 {
	t1All, t1Busy := getAllBusy(t1)
	t2All, t2Busy := getAllBusy(t2)

	if t2Busy <= t1Busy {
		return 0
	}
	if t2All <= t1All {
		return 100
	}
	return math.Min(100, math.Max(0, (t2Busy-t1Busy)/(t2All-t1All)*100))
}
func sharezoneCpuUtil( in []float64) float64 {
	length := float64(len(in))
	var total float64
	for _, v := range in {
		total += v
	}
	data := total/length
	return data
}
func main() {
	var configFile = flag.String("config", "/demo/42/config.conf", "Input your config file !")
	flag.Parse()
	config := InitConfig(*configFile)
	sharezoneCpuList := config["sharezoneCpuList"]
	cpuList := strings.Split(sharezoneCpuList," ")
	fmt.Println("cpuList=", cpuList)
	cpuTime1, err := getCpuTime()
	if err != nil {
		fmt.Println(err)
		os.Exit(4)
	}
	time.Sleep(1 * time.Second)
	cpuTime2, _ := getCpuTime()
	sharezoneCpuPercentList := []float64{}
	for k, v := range cpuTime1 {
		for _, vv := range cpuList {
			if v.CPU == vv {
				cpuPercent := calculateBusy(v, cpuTime2[k])
				sharezoneCpuPercentList = append(sharezoneCpuPercentList, cpuPercent)
			}
		}
		
	}
	sharezoneCpuPercent := sharezoneCpuUtil(sharezoneCpuPercentList)
	totalCpuUtil, err := getTotalCpuUsage()
	if err != nil {
		fmt.Println(err)
		os.Exit(5)
	}
	fmt.Println("host.cpu.util.sharezone: ", sharezoneCpuPercent)
	fmt.Println("host.cpu.util: ", totalCpuUtil)
}
