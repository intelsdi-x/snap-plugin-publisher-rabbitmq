// +build small

/*
http://www.apache.org/licenses/LICENSE-2.0.txt


Copyright 2015 Intel Corporation

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

package rmq

import (
	"testing"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/cpolicy"
	"github.com/intelsdi-x/snap/core/ctypes"
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
		configPolicy, err := rp.GetConfigPolicy()
		Convey("so rp.GetConfigPolicy should return a ConfigPolicy", func() {
			Convey("so config policy should not be nil", func() {
				So(configPolicy, ShouldNotBeNil)
			})
			Convey("so retreiving a config policy should not return an error", func() {
				So(err, ShouldBeNil)
			})
			Convey("so config policy should be a cpolicy.ConfigPolicy", func() {
				So(configPolicy, ShouldHaveSameTypeAs, &cpolicy.ConfigPolicy{})
			})

			Convey("so processing of configuration without optional parameter should return correct values of parameters", func() {
				testConfig := make(map[string]ctypes.ConfigValue)
				testConfig["uri"] = ctypes.ConfigValueStr{Value: "localhost:5672"}
				testConfig["exchange_name"] = ctypes.ConfigValueStr{Value: "snap"}
				testConfig["exchange_type"] = ctypes.ConfigValueStr{Value: "fanout"}
				testConfig["routing_key"] = ctypes.ConfigValueStr{Value: "metrics"}

				cfg, errs := configPolicy.Get([]string{""}).Process(testConfig)
				Convey("so configPolicy should process test config and return a config", func() {
					So(cfg, ShouldNotBeNil)
				})

				Convey("so parameters should have correct values", func() {
					So((*cfg)["uri"].(ctypes.ConfigValueStr).Value, ShouldEqual, "localhost:5672")
					So((*cfg)["exchange_name"].(ctypes.ConfigValueStr).Value, ShouldEqual, "snap")
					So((*cfg)["exchange_type"].(ctypes.ConfigValueStr).Value, ShouldEqual, "fanout")
					So((*cfg)["routing_key"].(ctypes.ConfigValueStr).Value, ShouldEqual, "metrics")
					So((*cfg)["durable"].(ctypes.ConfigValueBool).Value, ShouldEqual, true)
				})

				Convey("so testConfig processing should return no errors", func() {
					So(errs.HasErrors(), ShouldBeFalse)
				})
			})

			Convey("so processing of configuration with optional parameter should return correct values of parameters", func() {
				testConfig := make(map[string]ctypes.ConfigValue)
				testConfig["uri"] = ctypes.ConfigValueStr{Value: "localhost:5672"}
				testConfig["exchange_name"] = ctypes.ConfigValueStr{Value: "snap"}
				testConfig["exchange_type"] = ctypes.ConfigValueStr{Value: "fanout"}
				testConfig["routing_key"] = ctypes.ConfigValueStr{Value: "metrics"}
				testConfig["durable"] = ctypes.ConfigValueBool{Value: false}

				cfg, errs := configPolicy.Get([]string{""}).Process(testConfig)
				Convey("so configPolicy should process test config and return a config", func() {
					So(cfg, ShouldNotBeNil)
				})

				Convey("so parameters should have correct values", func() {
					So((*cfg)["uri"].(ctypes.ConfigValueStr).Value, ShouldEqual, "localhost:5672")
					So((*cfg)["exchange_name"].(ctypes.ConfigValueStr).Value, ShouldEqual, "snap")
					So((*cfg)["exchange_type"].(ctypes.ConfigValueStr).Value, ShouldEqual, "fanout")
					So((*cfg)["routing_key"].(ctypes.ConfigValueStr).Value, ShouldEqual, "metrics")
					So((*cfg)["durable"].(ctypes.ConfigValueBool).Value, ShouldEqual, false)
				})

				Convey("so testConfig processing should return no errors", func() {
					So(errs.HasErrors(), ShouldBeFalse)
				})
			})
		})
	})
}
