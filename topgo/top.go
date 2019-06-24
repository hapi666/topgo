package topgo

type TCPU struct {
	User       uint64
	Nice       uint64
	System     uint64
	Idle       uint64
	Iowait     uint64
	Irq        uint64
	Softirq    uint64
	Steal      uint64
	Guest      uint64
	Guest_nice uint64
	Total      uint64
}

type PCPU struct {
	Utime  uint64
	Stime  uint64
	Cutime uint64
	Cstime uint64
	Total  uint64
}

type PMem struct {
	Mem uint64
}

type UpTime struct {
	Time float64
}

type CPUs struct {
	UserAvg    float64
	NiceAvg    float64
	SystemAvg  float64
	IdelAvg    float64
	IowaitAvg  float64
	IrqAvg     float64
	SoftIrqAvg float64
	StealAvg   float64
	GuestAvg   float64
}

type LoadAvg struct {
	One     float64
	Five    float64
	Fifteen float64
}

type MemAndSwap struct {
	TotalMem     uint64
	FreeMem      uint64
	UsedMem      uint64
	AvailableMem uint64
	Buffer       uint64
	Cache        uint64
	TotalSwap    uint64
	FreeSwap     uint64
	UsedSwap     uint64
}

func GetPids() ([]int64, error) {
	return getPids()
}

func GetTotalCPUInfo() (c TCPU) {
	return getTotalCPUInfo()
}

func GetProcCPUInfo(pid int64) (c PCPU) {
	return getProcCPUInfo(pid)
}

func GetUpTime() string {
	return getUpTime()
}

func GetLoadAvg() (l LoadAvg) {
	return getLoadAvg()
}

func GetCPUs(first, second TCPU) (avg CPUs) {
	return getCPUs(first, second)
}

func GetMemAndSwap() (m MemAndSwap) {
	return getMemAndSwap()
}

func GetProcUid(pid int64) (uid int64) {
	return getProcUid(pid)
}

func GetProcUser(uid int64) (user string) {
	return getProcUser(uid)
}

func GetProcMemUsed(pid int64) (pm PMem) {
	return getProcMemUsed(pid)
}
