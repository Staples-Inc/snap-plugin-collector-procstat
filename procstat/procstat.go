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
	"strings"
	"time"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/cpolicy"
	"github.com/intelsdi-x/snap/core"
	"github.com/intelsdi-x/snap/core/ctypes"

	log "github.com/Sirupsen/logrus"
)

const (
	vendor        = "staples"
	fs            = "procfs"
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
	s := make(map[int32]Process)
	procstat := &Procstat{stats: s}
	return procstat
}

// Procstat defines procstat type
type Procstat struct {
	initialized bool
	pids        []Pid
	stats       map[int32]Process
}

func (p *Procstat) init(cfg map[string]ctypes.ConfigValue) error {
	if filesVal, ok := cfg["files"]; ok {
		files := strings.Split(filesVal.(ctypes.ConfigValueStr).Value, ",")
		for _, file := range files {
			pid, err := NewFilePid(file)
			if err != nil {
				log.Errorf("pid from filestring '%v' not parsed", file)
				continue
			}
			p.pids = append(p.pids, pid)
		}
	}
	p.initialized = true
	return nil
}

// registerPid if pid has been registered before, useful if process is restarted and pid changes
func (p *Procstat) registerPid(pid int32, processName string) (ps Process, err error) {
	ps, ok := p.stats[pid]
	if !ok {
		ps, err = NewProcess(pid, processName)
		if err != nil {
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

	for _, pid := range p.pids {
		pi, err := pid.GetPid()
		if err != nil {
			pid.MarkFailed()
			log.Errorf("failed to get pid from %v", pid.GetName())
		}
		processName := pid.GetName()

		// Check to see if pid has been registered before, useful if process is restarted or pid changes
		ps, err := p.registerPid(pi, processName)
		if err != nil {
			pid.MarkFailed()
			log.Errorf("error, unable to create process pid:%v name:%v", pi, processName)
			continue
		}

		fields, err := ps.GetStats()
		if err != nil {
			pid.MarkFailed()
			log.Errorf("error, pid:%v name:%v not collected, %v", pi, ps.GetName(), err)
			continue
		}

		for _, mt := range metricTypes {
			ns := make(core.Namespace, len(mt.Namespace()))
			copy(ns, mt.Namespace())
			var val interface{}
			var ok bool
			switch ns[3].Value {
			case processName, "*":
				if val, ok = fields[ns[4].Value]; !ok {
					continue
				}
			default:
				continue
			}
			ns[3].Value = ps.GetName()
			metric := plugin.MetricType{
				Namespace_: ns,
				Data_:      val,
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
		ns := core.NewNamespace(vendor, fs, pluginName).AddDynamicElement("processName", "Process Name").AddStaticElement(fields[name])
		appendMetric(&mts, ns, nil)
	}
	return mts, nil
}

func appendMetric(mts *[]plugin.MetricType, ns core.Namespace, data interface{}) {
	metric := plugin.MetricType{
		Namespace_: ns,
		Data_:      data,
		Timestamp_: time.Now(),
	}
	*mts = append(*mts, metric)
}

//GetConfigPolicy returns a ConfigPolicy
func (p *Procstat) GetConfigPolicy() (*cpolicy.ConfigPolicy, error) {
	cfg := cpolicy.New()
	pidFile, _ := cpolicy.NewStringRule("files", true)
	policy := cpolicy.NewPolicyNode()
	policy.Add(pidFile)
	cfg.Add([]string{vendor, fs, pluginName}, policy)
	return cfg, nil
}
