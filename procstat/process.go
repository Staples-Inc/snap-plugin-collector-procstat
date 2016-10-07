package procstat

import (
	"errors"
	"time"

	"github.com/shirou/gopsutil/process"
)

// Process acceses information through the /proc file system to return metrics on a process
type Process interface {
	GetStats() (map[string]interface{}, error)
	GetName() string
}

// Proc contains process information for collecting process metrics
type proc struct {
	ps   *process.Process
	name string
}

// NewProcess returns a new process object
func NewProcess(pid int32, processName string) (Process, error) {
	var prc proc
	p, err := process.NewProcess(pid)
	if err != nil {
		return &prc, errors.New("Unable to parse pid into process")
	}
	if processName == "" {
		processName, err = p.Name()
		if err != nil {
			return &prc, errors.New("Unable to get name from pid")
		}
	}
	prc = proc{ps: p, name: processName}

	return &prc, nil
}

// GetName returns the name of the process
func (p *proc) GetName() string {
	return p.name
}

// GetStats returns a map of proc stats
func (p *proc) GetStats() (map[string]interface{}, error) {
	fields := map[string]interface{}{}
	numThreads, err := p.ps.NumThreads()
	if err != nil {
		return fields, errors.New("unable to collect metrics, " + err.Error())
	}
	fds, err := p.ps.NumFDs()
	if err != nil {
		return fields, errors.New("unable to collect metrics, " + err.Error())
	}
	ctx, err := p.ps.NumCtxSwitches()
	if err != nil {
		return fields, errors.New("unable to collect metrics, " + err.Error())
	}
	io, err := p.ps.IOCounters()
	if err != nil {
		return fields, errors.New("unable to collect metrics, " + err.Error())
	}
	cpuTime, err := p.ps.Times()
	if err != nil {
		return fields, errors.New("unable to collect metrics, " + err.Error())
	}
	createTime, err := p.ps.CreateTime()
	if err != nil {
		return fields, errors.New("unable to collect metrics, " + err.Error())
	}
	cpuPerc, err := p.ps.Percent(0)
	if err != nil {
		return fields, errors.New("unable to collect metrics, " + err.Error())
	}
	mem, err := p.ps.MemoryInfo()
	if err != nil {
		return fields, errors.New("unable to collect metrics, " + err.Error())
	}

	fields["numThreads"] = numThreads
	fields["fds"] = fds
	fields["voluntary_context_switches"] = ctx.Voluntary
	fields["involuntary_context_switches"] = ctx.Involuntary
	fields["read_count"] = io.ReadCount
	fields["write_count"] = io.WriteCount
	fields["read_bytes"] = io.ReadBytes
	fields["write_bytes"] = io.WriteBytes
	fields["cpu_time_user"] = cpuTime.User
	fields["cpu_time_system"] = cpuTime.System
	fields["process_uptime"] = time.Now().Unix() - (createTime / 1000)
	fields["cpu_usage"] = cpuPerc
	fields["memory_rss"] = mem.RSS
	fields["memory_vms"] = mem.VMS
	fields["memory_swap"] = mem.Swap
	return fields, nil
}
