package procstat

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

// Pid queries for the pids configured for metrics collection
type Pid interface {
	GetPid() (int32, error)
	MarkFailed()
	GetName() string
}

type filePid struct {
	filepath string
	name     string
	pid      int32
	failed   bool
}

// NewFilePid returns a pid struct that made for reading pids from a pid file
func NewFilePid(filestring string) (Pid, error) {
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

func (f *filePid) MarkFailed() {
	f.failed = true
	return
}

func (f *filePid) GetName() string {
	return f.name
}

func (f *filePid) GetPid() (int32, error) {
	if f.failed {
		pid, err := getPidFromFile(f.filepath)
		if err != nil {
			return -1, errors.New("unable to read pid from file: " + f.filepath)
		}
		f.pid = pid
		f.failed = false
	}
	return f.pid, nil
}

func getPidFromFile(filename string) (int32, error) {
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
