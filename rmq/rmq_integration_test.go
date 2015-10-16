// +build integration

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
	"encoding/json"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/streadway/amqp"

	"github.com/intelsdi-x/pulse/control/plugin"
	"github.com/intelsdi-x/pulse/core/ctypes"
)

// integration test
func TestRmqIntegration(t *testing.T) {
	mt := plugin.PluginMetricType{
		Namespace_:          []string{"foo", "bar"},
		LastAdvertisedTime_: time.Now(),
		Version_:            1,
		Data_:               1,
	}
	data, _, err := plugin.MarshalPluginMetricTypes(plugin.PulseGOBContentType, []plugin.PluginMetricType{mt})
	Convey("Metric should encode successfully", t, func() {
		Convey("So err should be nil", func() {
			So(err, ShouldBeNil)
		})
	})
	rmqPub := NewRmqPublisher()
	cp, _ := rmqPub.GetConfigPolicy()
	config := map[string]ctypes.ConfigValue{
		"address":       ctypes.ConfigValueStr{Value: "127.0.0.1:5672"},
		"exchange_name": ctypes.ConfigValueStr{Value: "pulse"},
		"routing_key":   ctypes.ConfigValueStr{Value: "metrics"},
		"exchange_type": ctypes.ConfigValueStr{Value: "fanout"},
	}
	cfg, _ := cp.Get([]string{""}).Process(config)

	cKill := make(chan struct{})
	cMetrics, errc := connectToAmqp(cKill)
	err = rmqPub.Publish(plugin.PulseGOBContentType, data, *cfg)
	Convey("Publish should successfully publish metric to RabbitMQ server", t, func() {
		Convey("Publish data to RabbitMQ should not error", func() {
			So(err, ShouldBeNil)
		})
		Convey("We should be able to retrieve metric from RabbitMQ Server and validate", func() {
			Convey("Connecting to RabbitMQ server should not error", func() {
				So(errc, ShouldBeNil)
			})
			Convey("Validate metric", func() {
				if err == nil {
					select {
					case metric := <-cMetrics:
						var metrix []plugin.PluginMetricType
						err := json.Unmarshal(metric, &metrix)
						So(err, ShouldBeNil)
						So(metrix[0].Version(), ShouldEqual, mt.Version_)
						cKill <- struct{}{}
					case <-time.After(time.Second * 10):
						t.Fatal("Timeout when waiting for AMQP message")
					}

				}
			})
		})

	})
}

func connectToAmqp(cKill <-chan struct{}) (chan []byte, error) {
	conn, err := amqp.Dial("amqp://127.0.0.1:5672")
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	err = ch.ExchangeDeclare(
		"pulse",  //name
		"fanout", //kind
		true,     //durable
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return nil, err
	}

	q, err := ch.QueueDeclare(
		"",    // name
		true,  // durable
		false, // delete when usused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, err
	}
	//	FailOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(
		q.Name,    // queue name
		"metrics", // routing key
		"pulse",
		false,
		nil)
	if err != nil {
		return nil, err
	}
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)

	cMetrics := make(chan []byte)
	go func() {
		for {
			select {
			case msg := <-msgs:
				cMetrics <- msg.Body

			case _ = <-cKill:
				conn.Close()
				ch.Close()
				return
			}
		}
	}()
	return cMetrics, nil

}
