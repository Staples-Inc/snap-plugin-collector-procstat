// +build linux

/*
Copyright 2016 Staples, Inc.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
   http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package procstat

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/cpolicy"
	"github.com/intelsdi-x/snap/core"
	"github.com/intelsdi-x/snap/core/ctypes"

	log "github.com/Sirupsen/logrus"
	"github.com/shirou/gopsutil/process"
)

const (
	vendor        = "staples"
	pluginName    = "procstat"
	pluginType    = plugin.CollectorPluginType
	pluginVersion = 1
)

// Meta returns information about the plugin for the snap agent
func Meta() *plugin.PluginMeta {
	return plugin.NewPluginMeta(pluginName, pluginVersion, pluginType, []string{plugin.SnapGOBContentType}, []string{plugin.SnapGOBContentType})
}

// New returns a procstat plugin
func New() *Procstat {
	s := make(map[int32]*process.Process)
	procstat := &Procstat{stats: s}
	return procstat
}

// Procstat defines procstat type
type Procstat struct {
	initialized bool
	files       []string
	stats       map[int32]*process.Process
}

func (p *Procstat) init(cfg map[string]ctypes.ConfigValue) error {
	if filesVal, ok := cfg["files"]; ok {
		p.files = strings.Split(filesVal.(ctypes.ConfigValueStr).Value, ",")
	}
	p.initialized = true
	return nil
}

// checkPid if pid has been registered before, useful if process is restarted and pid changes
func(p *Procstat) checkPid(pid int32) (ps *process.Process, err error){
		ps, ok := p.stats[pid]
		if !ok {
			ps, err = process.NewProcess(pid)
			if err != nil {
				log.Errorf("pid %v from file unable to be accessed at /proc", pid)
				return nil, errors.New("Unable to access pid at /proc")
			}
			p.stats[pid] = ps
		}
		return ps, nil
}

// CollectMetrics returns metrics from gopsutil
func (p *Procstat) CollectMetrics(metricTypes []plugin.MetricType) ([]plugin.MetricType, error) {
	mts := []plugin.MetricType{}
	if !p.initialized {
		if err := p.init(metricTypes[0].Config().Table()); err != nil {
			return nil, err
		}
	}

	for _, pid := range p.files {
		pn := strings.Split(pid, ":")
		pidVal, err := getPidFromFile(pn[0])
		if err != nil {
			log.Errorf("%v, unable to read pid from file '%v'", err, pn[0])
			continue
		}

		// Check to see if pid has been registered before, useful if process is restarted
		ps, err := p.checkPid(pidVal)
		if err != nil {
			continue
		}


		// Check to see if process alias has been configured, if not use actual process name
		var processName string
		if len(pn) > 1 {
			processName = pn[1]
		} else {
			processName, err = ps.Name()
			if err != nil {
				log.Errorf("unable to parse process name of pid '%v' in '%v'", pidVal, pn[0])
				continue
			}
		}

		
		fields, err := p.getStats(ps)
		if err != nil {
			continue
		}
		for _, mt := range metricTypes {
			ns := make(core.Namespace, len(mt.Namespace()))
			copy(ns, mt.Namespace())
			var val interface{}
			var ok bool
			switch ns[3].Value {
				case processName, "*":
					if val, ok = fields[ns[4].Value]; !ok{
						continue
					}
				default:
					continue
			}
			ns[3].Value = processName
			metric := plugin.MetricType{
				Namespace_: ns,
				Data_:val,
				Timestamp_: time.Now(),
			}
			mts = append(mts, metric)
		}
	}
	p.initialized = true
	return mts, nil
}

// GetMetricTypes returns the metric types exposed by gopsutil
func (p *Procstat) GetMetricTypes(cfg plugin.ConfigType) (mts []plugin.MetricType, err error) {
	fields := []string{"numThreads", "fds", "voluntary_context_switches", "involuntary_context_switches", "read_count", "write_count", "read_bytes", "write_bytes", "cpu_time_user", "cpu_time_system", "process_uptime", "cpu_usage", "memory_rss", "memory_vms", "memory_swap"}
	for name := range fields {
		ns := core.NewNamespace(vendor, "procfs", pluginName).AddDynamicElement("processName", "Process Name").AddStaticElement(fields[name])
		appendMetric(&mts, ns, nil)
	}
	return mts, nil
}

func (p *Procstat) getStats(proc *process.Process) (map[string]interface{}, error) {
	fields := map[string]interface{}{}
	numThreads, err := proc.NumThreads()
	if err == nil {
		fields["numThreads"] = numThreads
	}

	fds, err := proc.NumFDs()
	if err == nil {
		fields["fds"] = fds
	}

	ctx, err := proc.NumCtxSwitches()
	if err == nil {
		fields["voluntary_context_switches"] = ctx.Voluntary
		fields["involuntary_context_switches"] = ctx.Involuntary
	}

	io, err := proc.IOCounters()
	if err == nil {
		fields["read_count"] = io.ReadCount
		fields["write_count"] = io.WriteCount
		fields["read_bytes"] = io.ReadBytes
		fields["write_bytes"] = io.WriteBytes
	}

	cpuTime, err := proc.Times()
	if err == nil {
		fields["cpu_time_user"] = cpuTime.User
		fields["cpu_time_system"] = cpuTime.System
	}

	createTime, err := proc.CreateTime()
	if err == nil {
		fields["process_uptime"] = time.Now().Unix() - (createTime / 1000)
	}

	cpuPerc, err := proc.Percent(0)
	if err == nil {
		fields["cpu_usage"] = cpuPerc
	}

	mem, err := proc.MemoryInfo()
	if err == nil {
		fields["memory_rss"] = mem.RSS
		fields["memory_vms"] = mem.VMS
		fields["memory_swap"] = mem.Swap
	}

	return fields, nil
}

func appendMetric(mts *[]plugin.MetricType, ns core.Namespace, data interface{}) {
	metric := plugin.MetricType{
		Namespace_: ns,
		Data_:      data,
		Timestamp_: time.Now(),
	}
	*mts = append(*mts, metric)
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

//GetConfigPolicy returns a ConfigPolicy
func (p *Procstat) GetConfigPolicy() (*cpolicy.ConfigPolicy, error) {
	cfg := cpolicy.New()
	pidFile, _ := cpolicy.NewStringRule("files", true)
	policy := cpolicy.NewPolicyNode()
	policy.Add(pidFile)
	cfg.Add([]string{vendor, "procfs", pluginName}, policy)
	return cfg, nil
}
