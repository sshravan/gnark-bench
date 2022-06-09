package cmd

import (
	"fmt"
	"strconv"

	. "github.com/klauspost/cpuid/v2"
)

// indicates if we compiled the binary with amd64_adx build tag
var amd64_adx bool

type benchData struct {
	Backend             string
	Curve               string
	Algorithm           string
	NbConstraints       int
	NbInternalVariables int
	NbSecretVariables   int
	NbPublicVariables   int
	NbCoefficients      int
	MaxRAM              uint64
	RunTime             int64
	Throughput          int
	Count               int

	// CPU info
	// NbPhysicalCores int
	// ThreadsPerCore  int
	// LogicalCores    int
	// CachelineBytes  int
	// L1DataBytes     int
	// L1InstrBytes    int
	// L2DataBytes     int
	// L3DataBytes     int
	// Frequency       int
	// SupportsADX     int
	// CPUName         string
}

func (bData benchData) headers() []string {
	return []string{"backend", "curve", "algorithm", "nbConstraints", "nbInternal", "nbSecret", "nbPublic", "nbCoefficients", "ram(MiB)", "time(ms)", "throughput(constraints/s)", "throughputPerCore(constraints/s)", "nbPhysicalCores", "nbThreadsPerCore", "nbLogicalCores", "cacheLine", "l1d", "l1i", "l2", "l3", "freq", "adx", "cpu", "amd64_adx", "count"}
}
func (bData benchData) values() []string {

	return []string{
		bData.Backend,
		bData.Curve,
		bData.Algorithm,
		strconv.Itoa(int(bData.NbConstraints)),
		strconv.Itoa(int(bData.NbInternalVariables)),
		strconv.Itoa(int(bData.NbSecretVariables)),
		strconv.Itoa(int(bData.NbPublicVariables)),
		strconv.Itoa(bData.NbCoefficients),
		strconv.Itoa(int(bData.MaxRAM)),
		strconv.Itoa(int(bData.RunTime)),
		strconv.Itoa(bData.Throughput),
		strconv.Itoa(bData.Throughput / CPU.LogicalCores),

		strconv.Itoa(CPU.PhysicalCores),
		strconv.Itoa(CPU.ThreadsPerCore),
		strconv.Itoa(CPU.LogicalCores),
		strconv.Itoa(CPU.CacheLine),
		strconv.Itoa(CPU.Cache.L1D),
		strconv.Itoa(CPU.Cache.L1I),
		strconv.Itoa(CPU.Cache.L2),
		strconv.Itoa(CPU.Cache.L3),
		strconv.Itoa(int(CPU.Hz / 1000000)),
		fmt.Sprintf("%v", CPU.Supports(ADX) && CPU.Supports(BMI2)),
		CPU.BrandName,
		fmt.Sprintf("%v", amd64_adx),
	}
}
