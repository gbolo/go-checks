// +build linux

package ps

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// UnixProcess is an implementation of Process that contains Unix-specific
// fields and information.
type UnixProcess struct {

	// http://man7.org/linux/man-pages/man5/proc.5.html

	// found in /proc/pid/stat
	pid   		int
	comm  		string
	state 		string
	ppid  		int
	pgrp  		int
	sid   		int
	num_threads 	int
	vsize 		int64
	rss 		int64

	// found in /proc/pid/cmdline
	cmdline 	string

	// found by stating /proc/pid directory
	uptime 		int64


}

func (p *UnixProcess) Pid() int {
	return p.pid
}

func (p *UnixProcess) PPid() int {
	return p.ppid
}

func (p *UnixProcess) Executable() string {
	return p.comm
}

func (p *UnixProcess) Uptime() int64 {
	return p.uptime
}

func (p *UnixProcess) Memory() int64 {
	return p.rss
}

// Refresh reloads all the data associated with this process.
func (p *UnixProcess) Refresh() error {

	statPath := fmt.Sprintf("/proc/%d/stat", p.pid)
	cmdLinePath := fmt.Sprintf("/proc/%d/cmdline", p.pid)

	fi, err := os.Stat(fmt.Sprintf("/proc/%d", p.pid))
	if err != nil {
		return err
	}

	// calculate process uptime by stat of pid dir
	stat := fi.Sys().(*syscall.Stat_t)
	p.uptime = time.Now().Unix() - int64(stat.Ctim.Sec)

	// read in required files
	dataBytes, err := ioutil.ReadFile(statPath)
	if err != nil {
		return err
	}

	dataBytesCmd, err := ioutil.ReadFile(cmdLinePath)
	if err != nil {
		return err
	}

	// for cmdline; we need to replace null bytes with spaces, except on last byte
	dataBytesCmdHr := []byte{}
	for k, b := range dataBytesCmd {
		if b == 0 && k != len(dataBytesCmd) - 1 {
			dataBytesCmdHr = append(dataBytesCmdHr, 040)
		} else {
			dataBytesCmdHr = append(dataBytesCmdHr, b)
		}
	}
	p.cmdline = string(dataBytesCmdHr)

	// parse the stat file
	data := string(dataBytes)
	psstat := strings.Split(data, " ")

	// populate remaining fields
	p.comm = psstat[1][1:len(psstat[1])-1]
	p.state = psstat[2]
	p.ppid, _ = strconv.Atoi(psstat[3])
	p.pgrp, _ = strconv.Atoi(psstat[4])
	p.sid, _ = strconv.Atoi(psstat[5])
	p.num_threads, _ = strconv.Atoi(psstat[19])
	p.vsize, _ = strconv.ParseInt(psstat[22], 10, 64)
	p.rss, _ = strconv.ParseInt(psstat[23], 10, 64)

	return err
}

func findProcess(pid int) (Process, error) {
	dir := fmt.Sprintf("/proc/%d", pid)
	_, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}

		return nil, err
	}

	return newUnixProcess(pid)
}

func processes() ([]Process, error) {
	d, err := os.Open("/proc")
	if err != nil {
		return nil, err
	}
	defer d.Close()

	results := make([]Process, 0, 50)
	for {
		fis, err := d.Readdir(10)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		for _, fi := range fis {
			// We only care about directories, since all pids are dirs
			if !fi.IsDir() {
				continue
			}

			// We only care if the name starts with a numeric
			name := fi.Name()
			if name[0] < '0' || name[0] > '9' {
				continue
			}

			// From this point forward, any errors we just ignore, because
			// it might simply be that the process doesn't exist anymore.
			pid, err := strconv.ParseInt(name, 10, 0)
			if err != nil {
				continue
			}

			p, err := newUnixProcess(int(pid))
			if err != nil {
				continue
			}

			results = append(results, p)
		}
	}

	return results, nil
}

func newUnixProcess(pid int) (*UnixProcess, error) {
	p := &UnixProcess{pid: pid}
	return p, p.Refresh()
}
