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
