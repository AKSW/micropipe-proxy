# Micropipe-proxy

> RabbitMQ to REST helper proxy for pipeline-based microservices

Micropipe proxy is a helper proxy for building pipeline-based microservices.
The proxy uses `config.yml` file to get all the required information about current service.
After reading the config file, the proxy will automatically connect to RabbitMQ and expose basic REST API for microservice to use.
It will also handle data validation, healthchecks and will send keep-alive messages to the bus (if enabled).

## Features

- Automatic connection, queue generation and listener setup for RabbitMQ
- Proxying of messages from and to RabbitMQ 
- Service versioning (via config file)
- Input and config schema validations
- Automated health monitoring and healthchecks
- Optional possibility to send heartbeats with service info to RabbitMQ

## Usage

We recommend using the proxy with Docker container.
You can find complete example service in `./example` folder.

1. Clone this repository and compile the binary for your platform (linux target is available as make target `make build-linux`)
2. Create config file describing your service
3. In your Dockerfile, copy proxy binary and config file into container
4. Set proxy as a container entry point 
5. Optionally add `HEALTHCHECK` docker command and point it to the proxy (e.g. `HEALTHCHECK --interval=5s --timeout=1s CMD curl -f http://localhost:8080/health || exit 1`)
6. Build and run your container while providing required environmental variables from Configuration section

Service route is built using ID and version fields from config file and looks as follows: `{ID}-{version}.#`.

## Configuration

Following things can be configured using environmental variables:

- `MICROPROXY_RABBIT_HOST` - address of RabbitMQ host to connect to (defaults to `localhost`)
- `MICROPROXY_EXCHANGE` - RabbitMQ exchange name to use (defaults to `microproxy`)
- `MICROPROXY_SERVER_LISTEN` - listen address for local REST API (defaults to `:8080`)
- `MICROPROXY_HEARTBEATS` - whether microproxy should send periodic heartbeats with service info into RabbitMQ
- `MICROPROXY_HEARTBEAT_ROUTE` - route to use for heartbeats (defaults to `microproxy.heartbeats`) 

## Config file

Configuration file is written in yaml format and contains the following fields:

- `id` - your service ID. Will be used during route generation 
- `name` - your service name. Only used as meta information for user
- `description` - your service description. Only used as meta information for user
- `version` - your service version. Will be used during route generation
- `command` - command to be executed as child service
- `responseEndpoint` - endpoint that proxy will use to deliver messages from RabbitMQ
- `inputSchema` - JSON schema describing the format of input messages. Used for incoming messages validation
- `outputSchema` - JSON schema describing the format of output messages. Only used as meta information for user
- `configSchema` - JSON schema describing the format of configuration fields in messages. Used for incoming messages validation

## TODO

- [ ] Unit tests
- [ ] Configure CI

## License

[MIT](https://opensource.org/licenses/MIT)
