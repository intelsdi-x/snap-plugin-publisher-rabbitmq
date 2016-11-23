# Snap publisher plugin - RabbitMQ

This plugin publishing metrics to AMQP queues via RabbitMQ

1. [Getting Started](#getting-started)
  * [System Requirements](#system-requirements)
  * [Installation](#installation)
  * [Configuration and Usage](#configuration-and-usage)
2. [Documentation](#documentation)
  * [Task manifest](#task-manifest)
  * [Examples](#examples)
  * [Roadmap](#roadmap)
3. [Community Support](#community-support)
4. [Contributing](#contributing)
5. [License](#license)
6. [Acknowledgements](#acknowledgements)

### System Requirements
* RabbitMQ set up and working on setup
* Plugin is working with AMQP protocol version 0.9.1
* Plugin supports only Linux systems

### Installation
#### Download RabbitMQ plugin binary:
You can get the pre-built binaries for your OS and architecture at plugin's [Github Releases](https://github.com/intelsdi-x/snap-plugin-publisher-rabbitmq/releases) page.

#### To build the plugin binary:
Fork https://github.com/intelsdi-x/snap-plugin-publisher-rabbitmq

Clone repo into `$GOPATH/src/github/intelsdi-x/`:
```
$ git clone https://github.com/<yourGithubID>/snap-plugin-publisher-rabbitmq
```
Build the plugin by running make in repo:
```
$ make
```
This builds the plugin in `./build`

### Configuration and Usage
* Set up the [snap framework](https://github.com/intelsdi-x/snap/blob/master/README.md#getting-started)

## Documentation
[rabbitmq](https://www.rabbitmq.com/documentation.html)

###  Task manifest
User need to provide following parameters in configuration for publisher:
- `uri` -  AMQP URI scheme for connecting to server. URI supports the following form: "user:password@ip:port/vhost".
           User, password, and vhost are optional. When not specified, the default values for RabbitMQ will be used.
- `exchange_name` - name of exchange,
- `routing_key` - routing key,
- `exchange_type` -  type of exchange,
- `durable` - sets durability (default: true).

### Examples

Example of running [psutil collector plugin](https://github.com/intelsdi-x/snap-plugin-collector-psutil) and publishing data to RabbitMQ.

Set up the [Snap framework](https://github.com/intelsdi-x/snap/blob/master/README.md#getting-started)

Ensure [Snap daemon is running](https://github.com/intelsdi-x/snap#running-snap):
* initd: `service snap-telemetry start`
* systemd: `systemctl start snap-telemetry`
* command line: `sudo snapteld -l 1 -t 0 &`

Download and load Snap plugins (paths to binary files for Linux/amd64):
```
$ wget http://snap.ci.snap-telemetry.io/plugins/snap-plugin-publisher-rabbitmq/latest/linux/x86_64/snap-plugin-publisher-rabbitmq
$ wget http://snap.ci.snap-telemetry.io/plugins/snap-plugin-collector-psutil/latest/linux/x86_64/snap-plugin-collector-psutil
$ snaptel plugin load snap-plugin-publisher-rabbitmq
$ snaptel plugin load snap-plugin-collector-psutil
```

Create a [task manifest](https://github.com/intelsdi-x/snap/blob/master/docs/TASKS.md) (see [exemplary tasks](examples/tasks/)),
for example `psutil-rabbitmq.json` with following content:
```json
{
  "version": 1,
  "schedule": {
    "type": "simple",
    "interval": "1s"
  },
  "workflow": {
    "collect": {
      "metrics": {
        "/intel/psutil/load/load1": {},
        "/intel/psutil/load/load15": {},
        "/intel/psutil/load/load5": {},
        "/intel/psutil/vm/available": {},
        "/intel/psutil/vm/free": {},
        "/intel/psutil/vm/used": {}
      },
      "publish": [
        {
          "plugin_name": "rabbitmq",
          "config": {
            "uri": "127.0.0.1:5672",
            "exchange_name": "snap",
            "routing_key": "metrics",
            "exchange_type": "fanout",
            "durable": true
          }
        }
      ]
    }
  }
}
```
Create a task:
```
$ snaptel task create -t psutil-rabbitmq.json
```

Watch created task:
```
$ snaptel task watch <task_id>
```

To stop previously created task:
```
$ snaptel task stop <task_id>
```

### Roadmap
There isn't a current roadmap for this plugin, but it is in active development. As we launch this plugin, we do not have any outstanding requirements for the next release.

If you have a feature request, please add it as an [issue](https://github.com/intelsdi-x/snap-plugin-publisher-rabbitmq/issues/new) and/or submit a [pull request](https://github.com/intelsdi-x/snap-plugin-publisher-rabbitmq/pulls).

## Community Support
This repository is one of **many** plugins in **Snap**, a powerful telemetry framework. See the full project at http://github.com/intelsdi-x/snap To reach out to other users, head to the [main framework](https://github.com/intelsdi-x/snap#community-support)

## Contributing
We love contributions!

There's more than one way to give back, from examples to blogs to code updates. See our recommended process in [CONTRIBUTING.md](CONTRIBUTING.md).

## License
[Snap](http://github.com:intelsdi-x/snap), along with this plugin, is an Open Source software released under the Apache 2.0 [License](LICENSE).

## Acknowledgements
List authors, co-authors and anyone you'd like to mention

* Author: [Nicholas Weaver](https://github.com/lynxbat)
* Author: [Szymon Konefal](https://github.com/skonefal)

And **thank you!** Your contribution, through code and participation, is incredibly important to us.
