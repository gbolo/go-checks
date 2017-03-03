package main

import (
	"fmt"
	"os"

	"github.com/gbolo/go-checks/lib/mem"
	"github.com/jessevdk/go-flags"
	"github.com/mackerelio/checkers"
)

var opts struct {
	UsedPct     int `short:"u" long:"used" value-name:"N" default:"90" description:"Used Memory Percent"`
	SwapUsedPct int `short:"s" long:"swapused" value-name:"N" default:"5" description:"Used Swap Percent"`
}

func main() {

	ckr := run(os.Args[1:])
	ckr.Name = "Memory"
	ckr.Exit()

}

func run(args []string) *checkers.Checker {

	_, err := flags.ParseArgs(&opts, args)
	if err != nil {
		os.Exit(1)
	}

	mem, _ := mem.Init()

	checkSt := checkers.OK
	msg := fmt.Sprintf("Used: %.2f%% Swap: %.2f%%",
		mem.UsedPercent(),
		mem.SwapUsedPercent(),
	)

	if float64(opts.UsedPct) < mem.UsedPercent() {
		checkSt = checkers.CRITICAL
	}

	if float64(opts.SwapUsedPct) < mem.SwapUsedPercent() {
		checkSt = checkers.CRITICAL
	}

	return checkers.NewChecker(checkSt, msg)

}
