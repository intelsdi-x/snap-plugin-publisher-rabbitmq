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
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/cpolicy"
	"github.com/intelsdi-x/snap/core/ctypes"
	"github.com/streadway/amqp"
)

func Meta() *plugin.PluginMeta {
	return plugin.NewPluginMeta(name, version, pluginType, []string{plugin.SnapGOBContentType}, []string{plugin.SnapGOBContentType})
}

type rmqPublisher struct{}

func NewRmqPublisher() *rmqPublisher {
	return &rmqPublisher{}

}

const (
	name       = "rabbitmq"
	version    = 8
	pluginType = plugin.PublisherPluginType
)

func (rmq *rmqPublisher) Publish(contentType string, content []byte, config map[string]ctypes.ConfigValue) error {
	logger := log.New()
	var metrics []plugin.MetricType
	switch contentType {
	case plugin.SnapGOBContentType:
		dec := gob.NewDecoder(bytes.NewBuffer(content))
		if err := dec.Decode(&metrics); err != nil {
			logger.Printf("Error decoding: error=%v content=%v", err, content)
			return err
		}
	default:
		logger.Printf("Error unknown content type '%v'", contentType)
		return fmt.Errorf("Unknown content type '%s'", contentType)
	}
	err := publishDataToRmq(
		metrics,
		config["address"].(ctypes.ConfigValueStr).Value,
		config["exchange_name"].(ctypes.ConfigValueStr).Value,
		config["routing_key"].(ctypes.ConfigValueStr).Value,
		config["exchange_type"].(ctypes.ConfigValueStr).Value,
		config["durable"].(ctypes.ConfigValueBool).Value,
		logger,
	)
	return err
}

func (rmq *rmqPublisher) GetConfigPolicy() (*cpolicy.ConfigPolicy, error) {
	cp := cpolicy.New()
	config := cpolicy.NewPolicyNode()

	r1, err := cpolicy.NewStringRule("address", true)
	handleErr(err)
	r1.Description = "RabbitMQ Address (host:port)"
	config.Add(r1)

	r2, err := cpolicy.NewStringRule("exchange_name", true)
	handleErr(err)
	r2.Description = "RabbitMQ Exchange Name"
	config.Add(r2)

	r3, err := cpolicy.NewStringRule("exchange_type", true)
	handleErr(err)
	r3.Description = "RabbitMQ Exchange Type"
	config.Add(r3)

	r4, err := cpolicy.NewStringRule("routing_key", true)
	handleErr(err)
	r4.Description = "RabbitMQ Routing Key"
	config.Add(r4)

	r5, err := cpolicy.NewBoolRule("durable", false, true)
	handleErr(err)
	r5.Description = "RabbitMQ Durable"
	config.Add(r5)

	cp.Add([]string{""}, config)
	return cp, nil
}

func publishDataToRmq(metrics []plugin.MetricType, address string, exName string, rtKey string, exKind string, durable bool, logger *log.Logger) error {
	conn, err := amqp.Dial("amqp://" + address)
	if err != nil {
		logger.Printf("RMQ Publisher: cannot open connection, %s", err)
		return err
	}
	defer conn.Close()

	c, err := conn.Channel()
	if err != nil {
		logger.Printf("RMQ Publisher: cannot open channel: %s", err)
		return err
	}

	err = c.ExchangeDeclare(exName, exKind, durable, false, false, false, nil)
	if err != nil {
		logger.Printf("RMQ Publisher: cannot declare exchange: %v", err)
		return err
	}

	data, err := json.Marshal(metrics)
	if err != nil {
		logger.Printf("RMQ Publisher: cannot marshal metrics: %v", err)
		return err
	}

	msg := amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		Timestamp:    time.Now(),
		ContentType:  "text/plain",
		Body:         data,
	}

	err = c.Publish(exName, rtKey, false, false, msg)
	if err != nil {
		logger.Printf("RMQ Publisher: cannot publish metric %v", err)
		return err
	}

	return nil
}

func handleErr(e error) {
	if e != nil {
		panic(e)
	}
}
