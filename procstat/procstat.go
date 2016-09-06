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
	procstat := &Procstat{}
	return procstat
}

// Procstat defines procstat type
type Procstat struct {
}

// CollectMetrics returns metrics from gopsutil
func (p *Procstat) CollectMetrics(metricTypes []plugin.MetricType) ([]plugin.MetricType, error) {
	pidsString := metricTypes[0].Config().Table()["files"]
	pids := strings.Split(pidsString.(ctypes.ConfigValueStr).Value, ",")
	mts := []plugin.MetricType{}
	for _, pid := range pids {
		pn := strings.Split(pid, ":")
		pid = pn[0]
		pidVal, err := getPidFromFile(pid)
		ps, err := process.NewProcess(pidVal)
		if err != nil {
			continue
		} else {
			fields, _ := getStats(ps)
			for _, metricType := range metricTypes {
				ns := metricType.Namespace()

				val, ok := getMapValueByNamespace(fields, ns[3:].Strings())
				if ok != nil {
					continue
				}
				appendMetric(&mts, ns, val)
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
		pid = pn[0]
		pidVal, err := getPidFromFile(pid)
		if err != nil {
			return mts, err
		}
		ps, err := process.NewProcess(pidVal)
		if err != nil {
			continue
		} else {
			var processName string
			if len(pn) > 1 {
				processName = pn[1]
			} else {
				processName, err = ps.Name()
				if err != nil {
					return nil, err
				}
			}
			fields, err := getStats(ps)
			if err != nil {
				return nil, err
			}
			for name := range fields {
				fmt.Println(processName)
				ns := core.NewNamespace(vendor, pluginName).AddStaticElement(processName).AddStaticElement(name)
				appendMetric(&mts, ns, nil)
			}
		}
	}
	return mts, nil
}

func getStats(proc *process.Process) (map[string]interface{}, error) {
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
		fields["cpu_time_idle"] = cpuTime.Idle
		fields["cpu_time_nice"] = cpuTime.Nice
		fields["cpu_time_iowait"] = cpuTime.Iowait
		fields["cpu_time_irq"] = cpuTime.Irq
		fields["cpu_time_soft_irq"] = cpuTime.Softirq
		fields["cpu_time_stolen"] = cpuTime.Steal
		fields["cpu_time_guest"] = cpuTime.Guest
		fields["cpu_time_guest_nice"] = cpuTime.GuestNice
	}

	createTime, err := proc.CreateTime()
	if err == nil {
		fields["process_uptime"] = time.Now().Unix() - (createTime / 1000)
	}

	cpuPerc, err := proc.Percent(time.Duration(0))
	if err == nil && cpuPerc != 0 {
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
