// Package system defines the domain model for the machine specifications
// gofetch displays and the contracts used to collect them.
package system

import "time"

// SystemInfo aggregates every specification gofetch displays.
//
//nolint:revive // system.Info is taken by the Collector contract; SystemInfo is the clearest name for the aggregate.
type SystemInfo struct {
	Host   Host
	CPU    CPU
	Memory Memory
	Disk   Disk
	Board  Board
	GPU    GPU
}

// Host holds operating system and machine identity information.
type Host struct {
	OS       string // e.g. "Windows 11 Pro (build 22631)"
	Hostname string
	User     string
	Uptime   time.Duration
}

// CPU holds processor model, core counts, clock and usage.
type CPU struct {
	Model         string
	PhysicalCores int
	LogicalCores  int
	MHz           float64
	UsagePercent  float64
}

// Memory holds RAM usage figures.
type Memory struct {
	UsedBytes   uint64
	TotalBytes  uint64
	UsedPercent float64
}

// Disk holds usage figures for the main disk.
type Disk struct {
	Mount       string // e.g. "C:"
	UsedBytes   uint64
	TotalBytes  uint64
	UsedPercent float64
}

// Board holds motherboard identification.
type Board struct {
	Manufacturer string
	Model        string
}

// GPU is a placeholder in v1: the renderer prints only the label.
// Model and usage data come in v1.1.
type GPU struct {
	Model string
}

func (h Host) apply(s *SystemInfo)   { s.Host = h }
func (c CPU) apply(s *SystemInfo)    { s.CPU = c }
func (m Memory) apply(s *SystemInfo) { s.Memory = m }
func (d Disk) apply(s *SystemInfo)   { s.Disk = d }
func (b Board) apply(s *SystemInfo)  { s.Board = b }
func (g GPU) apply(s *SystemInfo)    { s.GPU = g }
