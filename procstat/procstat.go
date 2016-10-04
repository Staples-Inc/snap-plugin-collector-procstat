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

var (
	errConfigError  = errors.New("Config read error: no pid file specified.")
	errPidFileError = errors.New("Pid read error: specified file not valid.")
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
	stats map[int32]*process.Process
}

// CollectMetrics returns metrics from gopsutil
func (p *Procstat) CollectMetrics(metricTypes []plugin.MetricType) ([]plugin.MetricType, error) {
	pidsString, ok := metricTypes[0].Config().Table()["files"]
	if !ok {
		return nil, log.Errorf("Unable to config file")
	}
	pids := strings.Split(pidsString.(ctypes.ConfigValueStr).Value, ",")
	mts := []plugin.MetricType{}
	for _, pid := range pids {
		pn := strings.Split(pid, ":")
		pidVal, err := getPidFromFile(pn[0])
		if err != nil {
			log.Errorf("%v, unable to read pid from file '%v'", err, pn[0])
			continue
		}

		ps, ok := p.stats[pidVal]
		if !ok {
			ps, err = process.NewProcess(pidVal)
			if err != nil {
				log.Errorf("pid %v from file '%v' unable to be accessed at /proc", pidVal, pn[0])
				continue
			}
			p.stats[pidVal] = ps
		}

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

		fields, _ := p.getStats(ps)
		for _, metricType := range metricTypes {
			ns := metricType.Namespace()
			if ns[3].Value == processName {
				val, ok := getMapValueByNamespace(fields, ns[4:].Strings())
				if ok != nil {
					continue
				}
				metricType.AddData(val)
				mts = append(mts, metricType)
			}
		}
	}
	return mts, nil
}

//getMapValueByNamespace gets value from map by namespace given in array of strings
func getMapValueByNamespace(m map[string]interface{}, ns []string) (val interface{}, err error) {
	if len(ns) == 0 {
		return val, fmt.Errorf("Namespace length equal to zero")
	}

	current := ns[0]

	var ok bool
	if len(ns) == 1 {
		if val, ok = m[current]; ok {
			return val, err
		}
		return val, fmt.Errorf("Key does not exist in map {key %s}", current)
	}

	if v, ok := m[current].(map[string]interface{}); ok {
		val, err = getMapValueByNamespace(v, ns[1:])
		return val, err
	}
	return val, err
}

// GetMetricTypes returns the metric types exposed by gopsutil
func (p *Procstat) GetMetricTypes(cfg plugin.ConfigType) ([]plugin.MetricType, error) {
	pidsString := cfg.Table()["files"]
	pids := strings.Split(pidsString.(ctypes.ConfigValueStr).Value, ",")
	mts := []plugin.MetricType{}
	for _, pid := range pids {
		pn := strings.Split(pid, ":")
		pidVal, err := getPidFromFile(pn[0])
		if err != nil {
			log.Errorf("unable to read pid from file '%v'", pn[0])
			continue
		}

		ps, ok := p.stats[pidVal]
		if !ok {
			ps, err = process.NewProcess(pidVal)
			if err != nil {
				log.Errorf("%v, pid %v from file '%v' unable to be accessed at /proc", err, pidVal, pn[0])
				continue
			}
		}

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
		for name := range fields {
			ns := core.NewNamespace(vendor, "procfs", pluginName).AddStaticElement(processName).AddStaticElement(name)
			appendMetric(&mts, ns, nil)
		}
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
	cfg.Add([]string{vendor, "procstat"}, policy)
	return cfg, nil
}
