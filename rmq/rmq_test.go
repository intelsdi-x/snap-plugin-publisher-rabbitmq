//
// +build unit

package rmq

import (
	"testing"

	"github.com/intelsdi-x/pulse/control/plugin"
	"github.com/intelsdi-x/pulse/control/plugin/cpolicy"
	"github.com/intelsdi-x/pulse/core/ctypes"
	. "github.com/smartystreets/goconvey/convey"
)

func TestRabbitmqPlugin(t *testing.T) {
	Convey("Meta should return Metadata for the plugin", t, func() {
		meta := Meta()
		So(meta.Name, ShouldResemble, name)
		So(meta.Version, ShouldResemble, version)
		So(meta.Type, ShouldResemble, plugin.PublisherPluginType)
	})

	Convey("Create RabbitMQ Publisher", t, func() {
		rp := NewRmqPublisher()
		Convey("so rp should not be nil", func() {
			So(rp, ShouldNotBeNil)
		})
		Convey("so rp should be of type rmqPublisher", func() {
			So(rp, ShouldHaveSameTypeAs, &rmqPublisher{})
		})
		Convey("so rp.GetConfigPolicy should return a ConfigPolicy", func() {
			configPolicy := rp.GetConfigPolicy()
			Convey("so config policy should not be nil", func() {
				So(configPolicy, ShouldNotBeNil)
			})
			Convey("so config policy should be a cpolicy.ConfigPolicy", func() {
				So(configPolicy, ShouldHaveSameTypeAs, cpolicy.ConfigPolicy{})
			})
			testConfig := make(map[string]ctypes.ConfigValue)
			testConfig["address"] = ctypes.ConfigValueStr{Value: "localhost:5672"}
			testConfig["exchange_name"] = ctypes.ConfigValueStr{Value: "pulse"}
			testConfig["exchange_type"] = ctypes.ConfigValueStr{Value: "fanout"}
			testConfig["routing_key"] = ctypes.ConfigValueStr{Value: "metrics"}
			cfg, errs := configPolicy.Get([]string{""}).Process(testConfig)
			Convey("so configPolicy should process test config and return a config", func() {
				So(cfg, ShouldNotBeNil)
			})
			Convey("so testConfig processing should return no errors", func() {
				So(errs.HasErrors(), ShouldBeFalse)
			})
		})
	})
}
