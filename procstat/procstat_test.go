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
	"testing"

	"github.com/intelsdi-x/snap/control/plugin"
	// "github.com/intelsdi-x/snap/control/plugin/cpolicy"
	// "github.com/intelsdi-x/snap/core"
	. "github.com/smartystreets/goconvey/convey"
)

type ProcInfoSuite struct {
	suite.Suite
	MockProcInfo string
}

func TestProcstatPlugin(t *testing.T) {
	meta := Meta()
	Convey("Meta should return metadata for the plugin", t, func() {
		Convey("So meta.Name should equal procstat", func() {
			So(meta.Name, ShouldEqual, pluginName)
		})
		Convey("So meta.Version should equal version", func() {
			So(meta.Version, ShouldEqual, pluginVersion)
		})
		Convey("So meta.Type should be of type plugin.ProcessorPluginType", func() {
			So(meta.Type, ShouldResemble, plugin.CollectorPluginType)
		})
	})
}

func mockNew() *Procstat {
	p := New()
	emptyCfg := make(map[string]ctypes.ConfigValue)
	p.init(emptyCfg)
	So(p, ShouldNotBeNil)
	So(p.snapMetricsNames, ShouldNotBeNil)
	return p
}

func (cis *ProcInfoSuite) TestGetMetricTypes() {
	_ = plugin.ConfigType{}
	Convey("Given cpu info plugin initialized", cis.T(), func() {
		p := mockNew()
		So(p, ShouldNotBeNil)
		Convey("When one wants to get list of available metrics", func() {
			mts, err := p.GetMetricTypes(plugin.ConfigType{})

			Convey("Then error should not be reported", func() {
				So(err, ShouldBeNil)
			})

			Convey("Then list of metrics is returned", func() {
				// Len mts = 24
				// cpuMetricsNumber = 3
				// len snapMetricsNames = 12
				So(len(mts), ShouldEqual, len(p.snapMetricsNames)*2)

				namespaces := []string{}
				for _, m := range mts {
					namespaces = append(namespaces, strings.Join(m.Namespace().Strings(), "/"))
				}

				So(namespaces, ShouldContain, "intel/procfs/cpu/*/user_percentage")
				So(namespaces, ShouldContain, "intel/procfs/cpu/*/nice_percentage")
				So(namespaces, ShouldContain, "intel/procfs/cpu/*/system_percentage")
				So(namespaces, ShouldContain, "intel/procfs/cpu/*/idle_percentage")
				So(namespaces, ShouldContain, "intel/procfs/cpu/*/iowait_percentage")
				So(namespaces, ShouldContain, "intel/procfs/cpu/*/irq_percentage")
				So(namespaces, ShouldContain, "intel/procfs/cpu/*/softirq_percentage")
				So(namespaces, ShouldContain, "intel/procfs/cpu/*/steal_percentage")
				So(namespaces, ShouldContain, "intel/procfs/cpu/*/guest_percentage")
				So(namespaces, ShouldContain, "intel/procfs/cpu/*/guest_nice_percentage")
				So(namespaces, ShouldContain, "intel/procfs/cpu/*/active_percentage")
				So(namespaces, ShouldContain, "intel/procfs/cpu/*/utilization_percentage")

				So(namespaces, ShouldContain, "intel/procfs/cpu/*/user_jiffies")
				So(namespaces, ShouldContain, "intel/procfs/cpu/*/nice_jiffies")
				So(namespaces, ShouldContain, "intel/procfs/cpu/*/system_jiffies")
				So(namespaces, ShouldContain, "intel/procfs/cpu/*/idle_jiffies")
				So(namespaces, ShouldContain, "intel/procfs/cpu/*/iowait_jiffies")
				So(namespaces, ShouldContain, "intel/procfs/cpu/*/irq_jiffies")
				So(namespaces, ShouldContain, "intel/procfs/cpu/*/softirq_jiffies")
				So(namespaces, ShouldContain, "intel/procfs/cpu/*/steal_jiffies")
				So(namespaces, ShouldContain, "intel/procfs/cpu/*/guest_jiffies")
				So(namespaces, ShouldContain, "intel/procfs/cpu/*/guest_nice_jiffies")
				So(namespaces, ShouldContain, "intel/procfs/cpu/*/active_jiffies")
				So(namespaces, ShouldContain, "intel/procfs/cpu/*/utilization_jiffies")

			})
		})
	})
}