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
