package system

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func getCPUUsage() (int64, int64, error) {
	file, err := os.Open("/proc/stat")
	if err != nil {
		return 0, 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	firstLine := scanner.Text()

	fields := strings.Fields(firstLine)
	if len(fields) < 5 {
		return 0, 0, fmt.Errorf("invalid /proc/stat format")
	}

	user, _ := strconv.ParseInt(fields[1], 10, 64)
	nice, _ := strconv.ParseInt(fields[2], 10, 64)
	system, _ := strconv.ParseInt(fields[3], 10, 64)
	idle, _ := strconv.ParseInt(fields[4], 10, 64)
	iowait, _ := strconv.ParseInt(fields[5], 10, 64)
	irq, _ := strconv.ParseInt(fields[6], 10, 64)
	softirq, _ := strconv.ParseInt(fields[7], 10, 64)

	idleAllTime := idle + iowait
	nonIdleAllTime := user + nice + system + irq + softirq
	totalTime := idleAllTime + nonIdleAllTime

	return totalTime, idleAllTime, nil
}

func MonitorCPU(doneChan <-chan struct{}, cpuChan chan<- int) {
	prevTotal, prevIdle, err := getCPUUsage()
	if err != nil {
		fmt.Println("Initial CPU reading error:", err)
		return
	}

	for {
		select {
		case <-doneChan:
			close(cpuChan)
			return
		default:
			time.Sleep(1 * time.Second)

			total, idle, err := getCPUUsage()
			if err != nil {
				fmt.Println("CPU reading error:", err)
				continue
			}

			totalDiff := total - prevTotal
			idleDiff := idle - prevIdle

			if totalDiff == 0 {
				continue
			}

			cpuUsage := int(100 * (totalDiff - idleDiff) / totalDiff)
			cpuChan <- cpuUsage

			prevTotal = total
			prevIdle = idle
		}
	}
}

func MonitorRAM(doneChan <-chan struct{}, ramChan chan<- int) {
	for {
		select {
		case <-doneChan:
			close(ramChan)
			return
		default:
			file, err := os.Open("/proc/meminfo")
			if err != nil {
				fmt.Println("Error opening /proc/meminfo:", err)
			}

			var total, available int64
			scanner := bufio.NewScanner(file)

			for scanner.Scan() {
				line := scanner.Text()
				if strings.HasPrefix(line, "MemTotal:") {
					fields := strings.Fields(line)
					total, _ = strconv.ParseInt(fields[1], 10, 64)
				} else if strings.HasPrefix(line, "MemAvailable:") {
					fields := strings.Fields(line)
					available, _ = strconv.ParseInt(fields[1], 10, 64)
				}
			}
			file.Close()

			usage := int((total - available) * 100 / total)
			ramChan <- usage

			time.Sleep(1 * time.Second)
		}
	}
}

func MonitorDisk(doneChan <-chan struct{}, diskChan chan<- int) {
	for {
		select {
		case <-doneChan:
			close(diskChan)
			return
		default:
			var stat syscall.Statfs_t
			err := syscall.Statfs("/", &stat)
			if err != nil {
				fmt.Println("Error getting disk stats:", err)
				continue
			}

			total := stat.Blocks * uint64(stat.Bsize)
			available := stat.Bfree * uint64(stat.Bsize)
			usage := int((total - available) * 100 / total)

			diskChan <- usage
			time.Sleep(1 * time.Minute)
		}
	}
}
