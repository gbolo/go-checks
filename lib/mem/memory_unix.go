// +build linux

package mem

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

const meminfo = "/proc/meminfo"

// fields and information.
type UnixMemory struct {

	// found in /proc/meminfo
	total     int64
	free      int64
	buffered  int64
	cached    int64
	available int64
	swapTotal int64
	swapFree  int64
}

func (m *UnixMemory) Total() int64 {
	return m.total
}

func (m *UnixMemory) Free() int64 {
	return m.free
}

func (m *UnixMemory) Buffered() int64 {
	return m.buffered
}

func (m *UnixMemory) Cached() int64 {
	return m.cached
}

func (m *UnixMemory) Available() int64 {
	return m.available
}

func (m *UnixMemory) UsedPercent() float64 {
	return float64(m.total-m.available) / float64(m.total) * 100
}

func (m *UnixMemory) SwapTotal() int64 {
	return m.swapTotal
}

func (m *UnixMemory) SwapFree() int64 {
	return m.swapFree
}

func (m *UnixMemory) SwapUsedPercent() float64 {
	return float64(m.swapTotal-m.swapFree) / float64(m.swapTotal) * 100
}

// Refresh reloads all the data associated with memory.
func (m *UnixMemory) Refresh() error {

	inFile, _ := os.Open(meminfo)
	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	for scanner.Scan() {
		l := scanner.Text()

		switch xs := strings.Split(l, ":"); xs[0] {
		case "MemTotal":
			m.total = extractInt(xs[1])
		case "MemFree":
			m.free = extractInt(xs[1])
		case "MemAvailable":
			m.available = extractInt(xs[1])
		case "Buffers":
			m.buffered = extractInt(xs[1])
		case "Cached":
			m.cached = extractInt(xs[1])
		case "SwapTotal":
			m.swapTotal = extractInt(xs[1])
		case "SwapFree":
			m.swapFree = extractInt(xs[1])
		}
	}

	return nil
}

func refreshMemory() (Memory, error) {
	m := &UnixMemory{}
	return m, m.Refresh()
}

// improve this function later
func extractInt(s string) int64 {

	// sample input:
	// MemTotal:        8048384 kB
	sInt := strings.Split(
		strings.TrimLeft(s, " "),
		" ",
	)[0]

	i, _ := strconv.ParseInt(sInt, 10, 64)

	return i

}
