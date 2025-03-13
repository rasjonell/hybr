package system

import (
	"bufio"
	"fmt"
	"github.com/rasjonell/hybr/internal/orchestration"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const (
	CPU_USAGE_EVENT  = orchestration.EventType("cpu_usage_event")
	RAM_USAGE_EVENT  = orchestration.EventType("ram_usage_event")
	DISK_USAGE_EVENT = orchestration.EventType("disk_usage_event")
)

type CPUUsageMonitor struct {
	EventType orchestration.EventType
}
type RAMUsageMonitor struct {
	EventType orchestration.EventType
}
type DiskUsageMonitor struct {
	EventType orchestration.EventType
}

func init() {
	manager := orchestration.GetSubscriptionManager()
	manager.RegisterEventSource(CPU_USAGE_EVENT, &CPUUsageMonitor{
		EventType: CPU_USAGE_EVENT,
	})
	manager.RegisterEventSource(RAM_USAGE_EVENT, &RAMUsageMonitor{
		EventType: RAM_USAGE_EVENT,
	})
	manager.RegisterEventSource(DISK_USAGE_EVENT, &DiskUsageMonitor{
		EventType: DISK_USAGE_EVENT,
	})
}

func (m *CPUUsageMonitor) Start(doneChan <-chan struct{}, cpuChan chan<- *orchestration.EventChannelData) {
	prevTotal, prevIdle, err := getCPUUsage()
	if err != nil {
		fmt.Println("Initial CPU reading error:", err)
		return
	}

	cpuUsage := int(100 * (prevTotal - prevIdle) / prevTotal)
	cpuChan <- orchestration.ToEventData(m.EventType, cpuUsage)

	for {
		select {
		case <-doneChan:
			return
		default:
			time.Sleep(10 * time.Second)

			total, idle, err := getCPUUsage()
			if err != nil {
				fmt.Println("CPU reading error:", err)
				continue
			}

			totalDiff := total - prevTotal
			idleDiff := idle - prevIdle

			cpuUsage := int(100 * (totalDiff - idleDiff) / totalDiff)
			cpuChan <- orchestration.ToEventData(m.EventType, cpuUsage)

			prevTotal = total
			prevIdle = idle
		}
	}
}

func (m *RAMUsageMonitor) Start(doneChan <-chan struct{}, ramChan chan<- *orchestration.EventChannelData) {
	for {
		select {
		case <-doneChan:
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
			ramChan <- orchestration.ToEventData(m.EventType, usage)

			time.Sleep(10 * time.Second)
		}
	}
}

func (m *DiskUsageMonitor) Start(doneChan <-chan struct{}, diskChan chan<- *orchestration.EventChannelData) {
	for {
		select {
		case <-doneChan:
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

			diskChan <- orchestration.ToEventData(m.EventType, usage)
			time.Sleep(10 * time.Minute)
		}
	}
}

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
