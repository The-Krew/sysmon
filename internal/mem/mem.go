package mem

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	"github.com/kopytkg/golog"
)

type Memory struct {
	Size     int
	Used     int
	Free     int
	Swap     int
	SwapFree int
	fields   map[string]int
}

func (m *Memory) getMemoryFields() {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		golog.Error(err.Error())
		os.Exit(1)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	m.fields = make(map[string]int)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")
		size, err := strconv.Atoi(strings.TrimSpace(strings.Fields(parts[1])[0]))
		if err != nil {
			golog.Error(err.Error())
			os.Exit(1)
		}
		m.fields[parts[0]] = size
	}
}

func (m *Memory) setMemoryFields() {
	m.Size = m.fields["MemTotal"]
	m.Used = m.fields["MemTotal"] - m.fields["MemAvailable"]
	m.Free = m.fields["MemAvailable"]
	m.Swap = m.fields["SwapTotal"]
	m.SwapFree = m.fields["SwapFree"]
}

func (m *Memory) Refresh() {
	m.getMemoryFields()
	m.setMemoryFields()
}

func NewMemory() *Memory {
	mem := Memory{}
	mem.getMemoryFields()
	mem.setMemoryFields()
	return &mem
}
