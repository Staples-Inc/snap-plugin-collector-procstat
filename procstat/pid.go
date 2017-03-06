package procstat

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strconv"
	"strings"
)

// Pid queries for the pids configured for metrics collection
type Pid interface {
	GetPids() ([]int32, error)
	MarkFailed()
	GetName() string
}

type filePid struct {
	filepath string
	name     string
	pid      []int32
	failed   bool
}

type patternPid struct {
	pattern string
	pid     []int32
	failed  bool
}

// NewFilePid returns a pid struct that made for reading pids from a pid file
func newFilePid(filestring string) (Pid, error) {
	ps := strings.Split(filestring, ":")
	// Check to see if process alias has been configured, if not use actual process name
	var pn string
	if len(ps) > 1 {
		pn = ps[1]
	} else {
		pn = ""
	}
	f := filePid{filepath: ps[0], name: pn, failed: true}
	return &f, nil
}

func newPgrepPid(pattern string) (Pid, error) {
	f := patternPid{pattern: pattern, failed: true}
	return &f, nil
}

func (f *filePid) MarkFailed() {
	f.failed = true
	return
}

func (f *filePid) GetName() string {
	return f.name
}

func (f *filePid) GetPids() ([]int32, error) {
	if f.failed {
		pid, err := discoverFile(f.filepath)
		if err != nil {
			return []int32{}, errors.New("unable to read pid from file: " + f.filepath)
		}
		f.pid = []int32{pid}
		f.failed = false
	}
	return f.pid, nil
}

func (f *patternPid) MarkFailed() {
	f.failed = true
	return
}

func (f *patternPid) GetName() string {
	return ""
}

func (f *patternPid) GetPids() ([]int32, error) {
	if f.failed {
		pids, err := discoverPgrep(f.pattern)
		if err != nil {
			return []int32{}, errors.New("unable to read pid from pgrep pattern: " + f.pattern)
		}
		f.pid = pids
		f.failed = false
	}
	return f.pid, nil
}

func discoverPgrep(pattern string) ([]int32, error) {
	var pids []int32
	bin, err := exec.LookPath("pgrep")
	if err != nil {
		return pids, fmt.Errorf("Couldn't find pgrep binary: %s", err)
	}
	pgrep, err := exec.Command(bin, "-f", pattern).Output()
	if err != nil {
		return pids, fmt.Errorf("Failed to execute %s. Error: '%s'", bin, err)
	}
	pgrepPids := strings.Fields(string(pgrep))
	for _, pid := range pgrepPids {
		ipid, err := strconv.Atoi(pid)
		if err == nil {
			pids = append(pids, int32(ipid))
		} else {
			return pids, err
		}
	}
	return pids, nil
}

func discoverFile(filename string) (int32, error) {
	var pid int
	pidString, err := ioutil.ReadFile(filename)
	if err != nil {
		return int32(pid), fmt.Errorf("Error: File not valid: %v", filename)
	}
	pid, err = strconv.Atoi(strings.TrimSpace(string(pidString[:])))
	if err != nil {
		return int32(pid), fmt.Errorf("Error: PID value could not be parsed from file %v", filename)
	}
	return int32(pid), nil
}
