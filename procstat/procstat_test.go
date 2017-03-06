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
	"strings"
	"testing"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/core"
	"github.com/intelsdi-x/snap/core/ctypes"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/suite"
)

const (
	proc_one = "psOne"
	proc_two = "psTwo"
)

type ProcInfoSuite struct {
	suite.Suite
	MockProcInfo string
}

type pidMock struct {
	name string
	pid  []int32
}

func newPidMock(n string, p int32) Pid {
	return &pidMock{name: n, pid: []int32{p}}
}

func (p *pidMock) GetPids() ([]int32, error) {
	return p.pid, nil
}
func (p *pidMock) MarkFailed() {
	return
}
func (p *pidMock) GetName() string {
	return p.name
}

type procMock struct {
	name string
}

func newProcessMock(processName string) Process {
	return &procMock{name: processName}
}

// GetName returns the name of the process
func (p *procMock) GetName() string {
	return p.name
}

// GetStats returns a map of proc stats
func (p *procMock) GetStats() (map[string]interface{}, error) {
	fields := make(map[string]interface{})
	fields["numThreads"] = 1
	fields["fds"] = 1
	fields["voluntary_context_switches"] = 1
	fields["involuntary_context_switches"] = 1
	fields["read_count"] = 1
	fields["write_count"] = 1
	fields["read_bytes"] = 1
	fields["write_bytes"] = 1
	fields["cpu_time_user"] = 1
	fields["cpu_time_system"] = 1
	fields["process_uptime"] = 1
	fields["cpu_usage"] = 1
	fields["memory_rss"] = 1
	fields["memory_vms"] = 1
	fields["memory_swap"] = 1
	return fields, nil
}

func mockNew() *Procstat {
	p := New()
	emptyCfg := make(map[string]ctypes.ConfigValue)
	p.stats[1] = newProcessMock(proc_one)
	p.stats[2] = newProcessMock(proc_two)
	p.rules = append(p.rules, newPidMock(proc_one, 1))
	p.rules = append(p.rules, newPidMock(proc_two, 2))
	p.init(emptyCfg)
	So(p, ShouldNotBeNil)
	p.initialized = true
	return p
}

func TestGetStatsSuite(t *testing.T) {
	suite.Run(t, &ProcInfoSuite{MockProcInfo: "ProcInfoSuite"})
}

func (cis *ProcInfoSuite) TestGetMetricTypes() {
	_ = plugin.ConfigType{}
	Convey("Given proc info plugin initialized", cis.T(), func() {
		p := mockNew()
		So(p, ShouldNotBeNil)
		Convey("When one wants to get list of available metrics", func() {
			mts, err := p.GetMetricTypes(plugin.ConfigType{})

			Convey("Then error should not be reported", func() {
				So(err, ShouldBeNil)
			})

			Convey("Then list of metrics is returned", func() {
				So(len(mts), ShouldEqual, 15)

				namespaces := []string{}
				for _, m := range mts {
					namespaces = append(namespaces, strings.Join(m.Namespace().Strings(), "/"))
				}
				So(namespaces, ShouldContain, "staples/procfs/procstat/*/cpu_time_system")
				So(namespaces, ShouldContain, "staples/procfs/procstat/*/cpu_time_user")
				So(namespaces, ShouldContain, "staples/procfs/procstat/*/cpu_usage")
				So(namespaces, ShouldContain, "staples/procfs/procstat/*/fds")
				So(namespaces, ShouldContain, "staples/procfs/procstat/*/involuntary_context_switches")
				So(namespaces, ShouldContain, "staples/procfs/procstat/*/memory_rss")
				So(namespaces, ShouldContain, "staples/procfs/procstat/*/memory_swap")
				So(namespaces, ShouldContain, "staples/procfs/procstat/*/memory_vms")
				So(namespaces, ShouldContain, "staples/procfs/procstat/*/numThreads")
				So(namespaces, ShouldContain, "staples/procfs/procstat/*/process_uptime")
				So(namespaces, ShouldContain, "staples/procfs/procstat/*/read_bytes")
				So(namespaces, ShouldContain, "staples/procfs/procstat/*/read_count")
				So(namespaces, ShouldContain, "staples/procfs/procstat/*/voluntary_context_switches")
				So(namespaces, ShouldContain, "staples/procfs/procstat/*/write_bytes")
				So(namespaces, ShouldContain, "staples/procfs/procstat/*/write_count")

			})
		})
	})
}

func (cis *ProcInfoSuite) TestCollectMetrics() {
	Convey("Given procstat plugin initialized", cis.T(), func() {
		p := mockNew()
		So(p, ShouldNotBeNil)

		Convey("When one wants to get values for given metric types", func() {
			processNames := []string{proc_one, proc_two}

			mTypes := []plugin.MetricType{}

			for _, processName := range processNames {
				mts := []plugin.MetricType{
					plugin.MetricType{Namespace_: core.NewNamespace(vendor, fs, pluginName, processName, "cpu_time_system")},
					plugin.MetricType{Namespace_: core.NewNamespace(vendor, fs, pluginName, processName, "cpu_time_user")},
					plugin.MetricType{Namespace_: core.NewNamespace(vendor, fs, pluginName, processName, "cpu_usage")},
					plugin.MetricType{Namespace_: core.NewNamespace(vendor, fs, pluginName, processName, "fds")},
					plugin.MetricType{Namespace_: core.NewNamespace(vendor, fs, pluginName, processName, "involuntary_context_switches")},
					plugin.MetricType{Namespace_: core.NewNamespace(vendor, fs, pluginName, processName, "memory_rss")},
					plugin.MetricType{Namespace_: core.NewNamespace(vendor, fs, pluginName, processName, "memory_swap")},
					plugin.MetricType{Namespace_: core.NewNamespace(vendor, fs, pluginName, processName, "memory_vms")},
					plugin.MetricType{Namespace_: core.NewNamespace(vendor, fs, pluginName, processName, "numThreads")},
					plugin.MetricType{Namespace_: core.NewNamespace(vendor, fs, pluginName, processName, "process_uptime")},
					plugin.MetricType{Namespace_: core.NewNamespace(vendor, fs, pluginName, processName, "read_bytes")},
					plugin.MetricType{Namespace_: core.NewNamespace(vendor, fs, pluginName, processName, "read_count")},
					plugin.MetricType{Namespace_: core.NewNamespace(vendor, fs, pluginName, processName, "voluntary_context_switches")},
					plugin.MetricType{Namespace_: core.NewNamespace(vendor, fs, pluginName, processName, "write_bytes")},
					plugin.MetricType{Namespace_: core.NewNamespace(vendor, fs, pluginName, processName, "write_count")},
				}
				mTypes = append(mTypes, mts...)
			}

			metrics, err := p.CollectMetrics(mTypes)

			Convey("Then no errors should be reported", func() {
				So(err, ShouldBeNil)
			})

			Convey("Then proper metrics values are returned", func() {
				So(len(metrics), ShouldEqual, len(mTypes))

				namespaces := []string{} //slice of namespaces for collected metrics

				for _, mt := range metrics {
					// Only jiffies metrics should not be nil
					for _, ns := range mt.Namespace_ {
						if strings.Contains(ns.Value, "jiffies") {
							So(mt.Data_, ShouldNotBeNil)
						}
					}
					//add namespace to slice of namespaces
					namespaces = append(namespaces, mt.Namespace().String())
				}
				Convey("Then collected metrics have desired namespaces", func() {
					for _, processName := range processNames {
						So(namespaces, ShouldContain, "/staples/procfs/procstat/"+processName+"/cpu_time_system")
						So(namespaces, ShouldContain, "/staples/procfs/procstat/"+processName+"/cpu_time_user")
						So(namespaces, ShouldContain, "/staples/procfs/procstat/"+processName+"/cpu_usage")
						So(namespaces, ShouldContain, "/staples/procfs/procstat/"+processName+"/fds")
						So(namespaces, ShouldContain, "/staples/procfs/procstat/"+processName+"/involuntary_context_switches")
						So(namespaces, ShouldContain, "/staples/procfs/procstat/"+processName+"/memory_rss")
						So(namespaces, ShouldContain, "/staples/procfs/procstat/"+processName+"/memory_swap")
						So(namespaces, ShouldContain, "/staples/procfs/procstat/"+processName+"/memory_vms")
						So(namespaces, ShouldContain, "/staples/procfs/procstat/"+processName+"/numThreads")
						So(namespaces, ShouldContain, "/staples/procfs/procstat/"+processName+"/process_uptime")
						So(namespaces, ShouldContain, "/staples/procfs/procstat/"+processName+"/read_bytes")
						So(namespaces, ShouldContain, "/staples/procfs/procstat/"+processName+"/read_count")
						So(namespaces, ShouldContain, "/staples/procfs/procstat/"+processName+"/voluntary_context_switches")
						So(namespaces, ShouldContain, "/staples/procfs/procstat/"+processName+"/write_bytes")
						So(namespaces, ShouldContain, "/staples/procfs/procstat/"+processName+"/write_count")
					}
				})
			})
		})
	})
}
