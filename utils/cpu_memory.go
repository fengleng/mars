package utils

import (
	"github.com/shirou/gopsutil/process"
	"os"
)

func GetCpuMemory() (cpu, memory float32) {
	var err error
	var p *process.Process
	p, err = process.NewProcess(int32(os.Getpid()))
	if err != nil {
		return
	}

	c, _ := p.CPUPercent()
	cpu = float32(c)

	m, _ := p.MemoryPercent()
	memory = m
	return
}
