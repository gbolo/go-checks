# Linux System Checks Written in Go
This repo contains a collection of checks written in go which can be statically compiled
and do not require any additional binaries in order to run.
These checks will get all needed information directly from `/proc`. 
This is useful for usage in **mininal docker containers** using alpine linux or just SCRATCH.

## Libraries
This repo contains the following libraries:

### github.com/gbolo/go-checks/lib/ps
Provides interface for running process information

### github.com/gbolo/go-checks/lib/mem
Provides interface for system memory information

## Binaries
This repo contains src to compile check binaries compatible with Sensu/Nagios

### github.com/gbolo/go-checks/bin/check-proc
```
$ check-proc --help
Usage:
  check-proc [OPTIONS]

Application Options:
  -p, --pid=N       The Process ID
  -u, --uptime=N    Required Uptime in seconds
  -r, --rss=N       Memory Usage in MB

Help Options:
  -h, --help        Show this help message
```

### github.com/gbolo/go-checks/bin/check-memory
```
$ check-memory --help
Usage:
  check-memory [OPTIONS]

Application Options:
  -u, --used=N        Used Memory Percent (default: 90)
  -s, --swapused=N    Used Swap Percent (default: 5)

Help Options:
  -h, --help          Show this help message
```