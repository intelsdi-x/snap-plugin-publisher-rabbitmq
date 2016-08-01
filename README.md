# snap publisher plugin - rabbitmq

This plugin publishing metrics to AMQP queues via RabbitMQ

1. [Getting Started](#getting-started)
  * [System Requirements](#system-requirements)
  * [Installation](#installation)
  * [Configuration and Usage](configuration-and-usage)
2. [Documentation](#documentation)
  * [Examples](#examples)
  * [Roadmap] (#roadmap)
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
You can get the pre-built binaries for your OS and architecture at snap's [Github Releases](https://github.com/intelsdi-x/snap/releases) page.

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
This builds the plugin in `/build/rootfs`

### Configuration and Usage
* Set up the [snap framework](https://github.com/intelsdi-x/snap/blob/master/README.md#getting-started)
* Ensure `$SNAP_PATH` is exported
`export SNAP_PATH=$GOPATH/src/github.com/intelsdi-x/snap/build`

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
Example task manifest to use RabbitMQ plugin:
```
{
    "version": 1,
    "schedule": {
        "type": "simple",
        "interval": "1s"
    },
    "workflow": {
        "collect": {
            "metrics": {
                "/intel/mock/foo": {},
                "/intel/mock/bar": {},
                "/intel/mock/*/baz": {}
            },
            "config": {
                "/intel/mock": {
                    "user": "root",
                    "password": "secret"
                }
            },
            "publish": [
                {
                    "plugin_name": "rabbitmq",
                    "config": {
                        "uri": "127.0.0.1:5672",
                        "exchange_name": "snap",
                        "routing_key": "metrics",
                        "exchange_type": "fanout",
                        "durable" : true
                    }
                }
            ]
        }
    }
}
```

### Roadmap
There isn't a current roadmap for this plugin, but it is in active development. As we launch this plugin, we do not have any outstanding requirements for the next release.

If you have a feature request, please add it as an [issue](https://github.com/intelsdi-x/snap-plugin-publisher-rabbitmq/issues/new) and/or submit a [pull request](https://github.com/intelsdi-x/snap-plugin-publisher-rabbitmq/pulls).

## Community Support
This repository is one of **many** plugins in **snap**, a powerful telemetry framework. See the full project at http://github.com/intelsdi-x/snap To reach out to other users, head to the [main framework](https://github.com/intelsdi-x/snap#community-support)

## Contributing
We love contributions!

There's more than one way to give back, from examples to blogs to code updates. See our recommended process in [CONTRIBUTING.md](CONTRIBUTING.md).

## License
[snap](http://github.com:intelsdi-x/snap), along with this plugin, is an Open Source software released under the Apache 2.0 [License](LICENSE).

## Acknowledgements
List authors, co-authors and anyone you'd like to mention

* Author: [Nicholas Weaver](https://github.com/lynxbat)
* Author: [Szymon Konefal](https://github.com/skonefal)

And **thank you!** Your contribution, through code and participation, is incredibly important to us.
