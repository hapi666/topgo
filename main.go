package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hapi666/topgo"
)

type stat struct {
	SysUpTime     string
	SysCurrCPU    topgo.TCPU
	SysLastCPU    topgo.TCPU
	SysCPUs       topgo.CPUss
	SysCPULoadAvg topgo.LoadAvg
	SysMemAndSwap topgo.MemAndSwap
	ProcCPU       topgo.PCPU
}

func getStat() *stat {
	s := stat{}
	s.SysCPULoadAvg = topgo.GetLoadAvg()
	s.SysUpTime = topgo.GetUpTime()
	s.SysCurrCPU = topgo.GetTotalCPUInfo()
	s.SysCPUs = topgo.GetCPUs(s.SysLastCPU, s.SysCurrCPU)
	s.SysLastCPU = s.SysCurrCPU
	s.SysMemAndSwap = topgo.GetMemAndSwap()
	return &s
}

func (s *stat) printTop() {
	fmt.Println("*************************************Top****************************************")
	fmt.Printf("go-top - %s  up %s,\t\tload average: %.2f, %.2f, %.2f\n",
		time.Now().Format("15:04"),
		s.SysUpTime,
		s.SysCPULoadAvg.
			One, s.SysCPULoadAvg.Five,
		s.SysCPULoadAvg.Fifteen)
	fmt.Printf("Cpu(s): %.1f%%us, %.1f%%sy, %.1f%%ni, %.1f%%id, %.1f%%wa, %.1f%%hi, %.1f%%si, %.1f%%st %.1f%%gu\n",
		s.SysCPUs.UserAvg,
		s.SysCPUs.SystemAvg,
		s.SysCPUs.NiceAvg,
		s.SysCPUs.IdelAvg,
		s.SysCPUs.IowaitAvg,
		s.SysCPUs.IrqAvg,
		s.SysCPUs.SoftIrqAvg,
		s.SysCPUs.StealAvg,
		s.SysCPUs.GuestAvg)
	fmt.Printf("Mem:  %9dk total, %9dk used, %9dk free, %9dk buffers\n",
		s.SysMemAndSwap.TotalMem,
		s.SysMemAndSwap.UsedMem,
		s.SysMemAndSwap.FreeMem,
		s.SysMemAndSwap.Buffer)
	fmt.Printf("Swap: %9dk total, %9dk used, %9dk free, %9dk cached\n",
		s.SysMemAndSwap.TotalSwap,
		s.SysMemAndSwap.UsedSwap,
		s.SysMemAndSwap.FreeSwap,
		s.SysMemAndSwap.Cache)

	fmt.Println("********************************************************************************")
	fmt.Println("PID\tUSER\t%CPU\t%MEM")

	//获取到目前运行的所有进程pid
	pids, err := topgo.GetPids()
	if err != nil {
		log.Printf("Failed to get pids: %v", err)
	}
	var onePcs []uint64
	var oneTcs []uint64
	//遍历每个进程
	for _, pid := range pids {
		//在第一个时间点求进程的CPU占用时间
		onePcs = append(onePcs, topgo.GetProcCPUInfo(pid).Total)
		oneTcs = append(oneTcs, topgo.GetTotalCPUInfo().Total)
	}
	//停1s
	time.Sleep(1 * time.Second)

	var twoPcs []uint64
	var twoTcs []uint64
	for _, pid := range pids {
		//在第二个时间点
		twoPcs = append(twoPcs, topgo.GetProcCPUInfo(pid).Total)
		twoTcs = append(twoTcs, topgo.GetTotalCPUInfo().Total)
	}
	for i, pid := range pids {
		memStr := strconv.FormatUint(topgo.GetProcMemUsed(pid).Mem/s.SysMemAndSwap.TotalMem, 10)
		f, err := strconv.ParseFloat(memStr, 64)
		if err != nil {
			log.Printf("Failed to parse memory utilization: %v", err)
		}
		fmt.Printf("%d\t%s\t%.2f\t%.2f\n",
			pids[i],
			topgo.GetProcUser(topgo.GetProcUid(pid)),
			float64(100*(twoPcs[i]-onePcs[i])/(twoTcs[i]-oneTcs[i])),
			f)
	}
}

func main() {
	s := getStat()
	s.printTop()
}
