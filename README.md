# snap collector plugin - procstat
The procstat collector takes a specified pid file and collects metrics from the process specified by that pid file.

[![Build Status](https://travis-ci.org/Staples-Inc/snap-plugin-collector-procstat.svg?branch=master)](https://travis-ci.org/Staples-Inc/snap-plugin-collector-procstat)

## Getting Started
### System Requirements
* [golang 1.5+](https://golang.org/dl/) - needed only for building
* [glide](https://github.com/Masterminds/glide) - needed for dependencies

### Operating systems
All OSs currently supported by plugin:
* Linux/amd64

### Build
Build the plugin by running make within the repo:
```
$ make
```
This builds the plugin in `/build/rootfs/`

### Run
The configuration for the procstat plugin is in the task file. The pids designated for collection must exist in a pid file.
```
"workflow": {
  "collect": {
    "metrics": {
      "/staples/procfs/procstat/*/cpu_time_system": {},
      "/staples/procfs/procstat/*/cpu_time_user": {},
      "/staples/procfs/procstat/*/cpu_usage": {},
      "/staples/procfs/procstat/*/fds": {},
      "/staples/procfs/procstat/*/involuntary_context_switches": {},
      "/staples/procfs/procstat/*/memory_rss": {},
      "/staples/procfs/procstat/*/memory_swap": {},
      "/staples/procfs/procstat/*/memory_vms": {},
      "/staples/procfs/procstat/*/numThreads": {},
      "/staples/procfs/procstat/*/process_uptime": {},
      "/staples/procfs/procstat/*/read_bytes": {},
      "/staples/procfs/procstat/*/read_count": {},
      "/staples/procfs/procstat/*/voluntary_context_switches": {},
      "/staples/procfs/procstat/*/write_bytes": {},
      "/staples/procfs/procstat/*/write_count": {}
    },
    "config": {
      "/staples/procfs/procstat":{
        "files":"/tmp/snap.pid:coolio,/tmp/syslog.pid"
      }
    },
    "publish": null
  }
}
```
* Specify the pid file with the desired name identifier using the format: `"files": "<filepath>:<name>,<filepath>:<name>,..."`
* Ensure your snap agent is run with the correct permissions to collect from the `/proc/<pid>/` file

## Roadmap
* Add pgrep functionality for tracking a process in addition to files.
```
"config": {
  "/staples/procfs/procstat":{
    "files":"/tmp/snap.pid:coolio,/tmp/syslog.pid",
    "pgrep":"mongo,java:cassandra"
  }
}
```
* Allow for pid interface to manage multiple pids if a file contains more than one or if pgrep returns more than one

If you have suggestions please open up an issue or provide a pull request.
