package topgo

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"
)

func getPids() ([]int64, error) {
	files, err := ioutil.ReadDir("/proc")
	if err != nil {
		log.Printf("Failed to open /proc directory: %v\n", err)
		return nil, err
	}
	pids := make([]int64, 0, 100)
	for _, file := range files {
		if file.IsDir() {
			pid, err := strconv.ParseInt(file.Name(), 0, 64)
			if err != nil {
				continue
			} else {
				pids = append(pids, pid)
			}
		}
	}
	return pids, nil
}
func getUpTime() string {
	info, err := ioutil.ReadFile("/proc/uptime")
	if err != nil {
		log.Printf("Failed to read /proc/uptime file: %v\n", err)
	}
	var u UpTime
	u.Time, err = strconv.ParseFloat(strings.Split(string(info), " ")[0], 64)
	if err != nil {
		log.Printf("topgo.getUpTime(): Error ParseFloat(info): %v\n", err)
	}

	t, err := time.ParseDuration(strconv.FormatFloat(u.Time, 'f', 5, 64) + "s")
	if err != nil {
		log.Printf("topgo.getUpTime(): Error ParseDuration(u.Time): %v\n", err)
	}
	days := t.Hours() / 24
	if days >= 1 && days < 365 {
		return strconv.FormatFloat(days, 'f', 0, 64) + "days"
	} else if days < 1 {
		return strconv.FormatFloat(t.Hours(), 'f', 0, 64) + "hours"
	} else { //days>365
		years := days / 365
		return strconv.FormatFloat(years, 'f', 0, 64) + "years"
	}
}

func getProcCPUInfo(pid int64) (c PCPU) {
	pidStr := strconv.FormatInt(pid, 10)
	info, err := ioutil.ReadFile("/proc" + "/" + pidStr + "/stat")
	if err != nil {
		log.Printf("Failed to read /proc/"+pidStr+"/stat file: %v\n", err)
	}
	content := strings.Split(string(info), " ")
	if len(content) > 0 {
		c.Utime, err = strconv.ParseUint(content[13], 10, 64)
		if err != nil {
			log.Printf("topgo.getProcCPUInfo():Error ParseUint(content[13]): %v\n", err)
			return
		}
		c.Stime, err = strconv.ParseUint(content[14], 10, 64)
		if err != nil {
			log.Printf("topgo.getProcCPUInfo():Error ParseUint(content[14]): %v\n", err)
			return
		}
		c.Cutime, err = strconv.ParseUint(content[15], 10, 64)
		if err != nil {
			log.Printf("topgo.getProcCPUInfo():Error ParseUint(content[15]): %v\n", err)
			return
		}
		c.Cstime, err = strconv.ParseUint(content[16], 10, 64)
		if err != nil {
			log.Printf("topgo.getProcCPUInfo():Error ParseUint(content[16]): %v\n", err)
			return
		}
	}
	c.Total = c.Stime + c.Utime + c.Cstime + c.Cutime
	return
}

func getTotalCPUInfo() (c TCPU) {
	info, err := ioutil.ReadFile("/proc/stat")
	if err != nil {
		log.Printf("Failed to read /proc/stat file: %v\n", err)
		return
	}
	line, _, err := bufio.NewReader(bytes.NewBuffer(info)).ReadLine()
	fields := strings.Fields(string(line))
	if len(fields) > 0 && fields[0] == "cpu" {
		parseTCPU(fields, &c)
	}
	c.Total = c.User + c.Nice + c.System + c.Idle + c.Iowait + c.Irq + c.Steal + c.Softirq + c.Guest
	return
}

func parseTCPU(fields []string, tcpu *TCPU) {
	for i, field := range fields {
		f, _ := strconv.ParseUint(field[1:], 10, 64)
		switch i {
		case 1:
			tcpu.User = f
		case 2:
			tcpu.Nice = f
		case 3:
			tcpu.System = f
		case 4:
			tcpu.Idle = f
		case 5:
			tcpu.Iowait = f
		case 6:
			tcpu.Irq = f
		case 7:
			tcpu.Softirq = f
		case 8:
			tcpu.Steal = f
		case 9:
			tcpu.Guest = f
		}
	}
}

func subtract(second, first uint64) float64 {
	return float64(second - first)
}

func getCPUs(first, second TCPU) (avg CPUs) {
	comm := 100 / subtract(second.Total, first.Total)
	avg.UserAvg = subtract(second.User, first.User) * comm
	avg.NiceAvg = subtract(second.Nice, first.Nice) * comm
	avg.SystemAvg = subtract(second.System, first.System) * comm
	avg.IdelAvg = subtract(second.Idle, first.Idle) * comm
	avg.IowaitAvg = subtract(second.Iowait, first.Iowait) * comm
	avg.IrqAvg = subtract(second.Irq, first.Irq) * comm
	avg.SoftIrqAvg = subtract(second.Softirq, first.Softirq) * comm
	avg.StealAvg = subtract(second.Steal, first.Steal) * comm
	avg.GuestAvg = subtract(second.Guest, first.Guest) * comm
	return
}

func getLoadAvg() (l LoadAvg) {
	info, err := ioutil.ReadFile("/proc/loadavg")
	if err != nil {
		log.Printf("Failed to read /proc/loadavg file: %v\n", err)
		return
	}
	line, _, err := bufio.NewReader(bytes.NewBuffer(info)).ReadLine()
	if err != nil {
		log.Printf("topgo.getLoadAvg(): Error ReadLine(): %v\n", err)
		return
	}
	fields := strings.Fields(string(line))
	for i, field := range fields {
		var temp float64
		if i < 3 {
			temp, err = strconv.ParseFloat(field, 64)
			if err != nil {
				log.Printf("topgo.getLoadAvg(): Error ParseFloat(%s): %v\n", field, err)
				return
			}
		}
		switch i {
		case 0:
			l.One = temp
		case 1:
			l.Five = temp
		case 2:
			l.Fifteen = temp
		}
	}
	return
}

func getMemAndSwap() (m MemAndSwap) {
	info, err := ioutil.ReadFile("/proc/meminfo")
	if err != nil {
		log.Printf("Failed to read /proc/meminfo: %v\n", err)
		return
	}
	reader := bufio.NewReader(bytes.NewBuffer(info))
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		fields := strings.Fields(string(line))
		comm, err := strconv.ParseUint(fields[1], 10, 64)
		if err != nil {
			log.Printf("topgo.getMemAndSwap(): Error ParseUint(%s): %v\n", fields[1], err)
			return
		}
		switch fields[0] {
		case "MemTotal:":
			m.TotalMem = comm
		case "MemFree:":
			m.FreeMem = comm
		case "MemAvailable:":
			m.AvailableMem = comm
		case "Buffers:":
			m.Buffer = comm
		case "Cached:":
			m.Cache = comm
		case "SwapTotal:":
			m.TotalSwap = comm
		case "SwapFree:":
			m.FreeSwap = comm
		}
	}
	m.UsedMem = m.TotalMem - m.FreeMem
	m.UsedSwap = m.TotalSwap - m.FreeSwap
	return
}

func getProcUid(pid int64) (uid int64) {
	pidStr := strconv.FormatInt(pid, 10)
	info, err := ioutil.ReadFile("/proc/" + pidStr + "/status")
	if err != nil {
		log.Printf("Failed to read /proc/"+pidStr+"/status file"+": %v\n", err)
		return
	}
	reader := bufio.NewReader(bytes.NewBuffer(info))
	for {
		content, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		fields := strings.Fields(string(content))
		if len(fields) > 0 && fields[0] == "Uid" {
			uid, err = strconv.ParseInt(fields[1], 10, 64)
			if err != nil {
				log.Printf("topgo.getProcUid(): Error ParseInt(%s): %v\n", fields[1], err)
				return
			}
			break
		}
	}
	return
}

func getProcUser(uid int64) (user string) {
	info, err := ioutil.ReadFile("/etc/passwd")
	if err != nil {
		log.Printf("Failed to read /etc/passwd file: %v\n", err)
		return
	}
	reader := bufio.NewReader(bytes.NewBuffer(info))
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		fields := strings.Split(string(line), ":")
		f, err := strconv.ParseInt(fields[2], 10, 64)
		if err != nil {
			log.Printf("topgo.getProcUser():Error ParseInt(%s): %v\n", fields[2], err)
		}
		if f == uid {
			user = fields[0]
		}
	}
	return
}

func getProcMemUsed(uid int64) (pm PMem) {
	uidStr := strconv.FormatInt(uid, 10)
	info, err := ioutil.ReadFile("/proc/" + uidStr + "/statm")
	if err != nil {
		log.Printf("Failed to read /proc/"+uidStr+"/statm"+": %v\n", err)
		return
	}
	content, _, err := bufio.NewReader(bytes.NewBuffer(info)).ReadLine()
	if err != nil {
		log.Printf("topgo.getProcMemUsed(): Error ReadLine(): %v\n", err)
		return
	}
	fields := strings.Fields(string(content))
	if len(fields) > 0 {
		pm.Mem, err = strconv.ParseUint(fields[1], 10, 64)
		if err != nil {
			log.Printf("topgo.getProcMemUsed(): Error ParseUint(%s): %v\n", fields[1], err)
			return
		}
	}
	return
}
