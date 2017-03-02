package main

import (
	"github.com/gbolo/go-checks/lib/ps"
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/mackerelio/checkers"
	"os"
)


var opts struct {
	Pid	      int     `short:"p" long:"pid" value-name:"N" description:"The Process ID"`
	Uptime        int64   `short:"u" long:"uptime" value-name:"N" description:"Required Uptime in seconds"`
	Rss           int64   `short:"r" long:"rss" value-name:"N" description:"Memory Usage in MB"`
}

func main(){

	ckr := run(os.Args[1:])
	ckr.Name = "Process"
	ckr.Exit()

}

func run(args []string) *checkers.Checker {

	_, err := flags.ParseArgs(&opts, args)
	if err != nil {
		os.Exit(1)
	}

	// default run with no pid specified
	if opts.Pid < 1 {
		procs, _ := ps.Processes()
		return checkers.NewChecker(
			checkers.OK,
			fmt.Sprintf("Found %d running processes", len(procs) ),
		)
	}

	proc, err := ps.FindProcess(opts.Pid)
	if proc == nil {
		return checkers.Critical(fmt.Sprintf("Process with ID %d is not running", opts.Pid ))
	}

	memoryRss := proc.Memory() * 4 / 1024

	checkSt := checkers.OK
	msg := fmt.Sprintf("%s running [%d MB] [%d secs]",
		proc.Executable(),
		memoryRss,
		proc.Uptime(),
	)

	// check proc uptime
	if opts.Uptime > 0 && opts.Uptime > proc.Uptime() {
		checkSt = checkers.CRITICAL
		msg = fmt.Sprintf("%s (%d) has not been running long enough: %d secs",
			proc.Executable(),
			proc.Pid(),
			proc.Uptime(),
		)
	}

	// check proc memory
	if opts.Rss > 0 && opts.Rss < memoryRss {
		checkSt = checkers.CRITICAL
		msg = fmt.Sprintf("%s (%d) is using too much memory: %d MB",
			proc.Executable(),
			proc.Pid(),
			memoryRss,
		)
	}

	return checkers.NewChecker(checkSt, msg)
}
