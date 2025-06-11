package cpu

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	"github.com/kopytkg/golog"
)

type Core map[string]any

type CPU struct {
	Cores       []Core
	Freq        []float64
	ModelName   string
	startIdles  []uint64
	startTotals []uint64
	Usage       []float64
}

func drillCPUInfo() []Core {
	file, err := os.Open("/proc/cpuinfo")
	if err != nil {
		golog.Error(err.Error())
		os.Exit(1)
	}
	defer file.Close()

	var cores []Core
	core := make(Core)

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			if len(core) > 0 {
				cores = append(cores, core)
				core = make(Core)
			}
			continue
		}

		parts := strings.SplitN(line, ":", 2)

		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])

			nval, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 32)
			if err != nil {
				core[key] = strings.TrimSpace(parts[1])
			} else {
				core[key] = nval
			}

		}
	}
	if len(core) > 0 {
		cores = append(cores, core)
	}

	return cores
}

func getCPUTimes() ([]uint64, []uint64) {
	file, err := os.Open("/proc/stat")
	if err != nil {
		golog.Error(err.Error())
		os.Exit(1)
	}

	defer file.Close()

	var idles []uint64
	var totals []uint64

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		var sum uint64 = 0
		if strings.HasPrefix(line, "cpu") {
			parts := strings.Fields(line)
			for i, part := range parts[1:] {

				nval, err := strconv.ParseUint(part, 10, 64)

				if err != nil {
					golog.Error(err.Error())
					break
				}
				sum += nval
				if i == 3 {
					idles = append(idles, uint64(nval))
				}

			}
			totals = append(totals, uint64(sum))
		}
	}
	return idles, totals
}

func NewCPU() *CPU {
	cpu := CPU{}

	cpu.Cores = drillCPUInfo()
	cpu.startIdles, cpu.startTotals = getCPUTimes()

	if model, ok := cpu.Cores[0]["model name"].(string); ok {
		cpu.ModelName = model
	}

	for _, core := range cpu.Cores {
		if freq, ok := core["cpu MHz"].(float64); ok {
			cpu.Freq = append(cpu.Freq, freq)
		}
	}

	return &cpu
}

func (c *CPU) Refresh() {
	for i, core := range drillCPUInfo() {
		if freq, ok := core["cpu MHz"].(float64); ok {
			c.Freq[i] = freq
		}
	}
}

func (c *CPU) UsageRefresh() {
	idles, totals := getCPUTimes()
	if len(c.Usage) != len(idles) {
		c.Usage = make([]float64, len(idles))
	}
	for i := 0; i < len(idles); i++ {
		idleDelta := idles[i] - c.startIdles[i]
		totalDelta := totals[i] - c.startTotals[i]
		if totalDelta == 0 {
			c.Usage[i] = 0
			continue
		}
		c.Usage[i] = 100.0 * (1.0 - float64(idleDelta)/float64(totalDelta))
	}
	c.startIdles = idles
	c.startTotals = totals
}
