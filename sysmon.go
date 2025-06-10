package sysmon

import (
	"github.com/the-krew/sysmon/internal/cpu"
	"github.com/the-krew/sysmon/internal/mem"
)

type SysmonAPI struct {
	CPU      cpu.CPU
	RAM      mem.Memory
	OnUpdate func()
}
