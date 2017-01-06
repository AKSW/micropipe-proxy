# FOX service wrapper example

Contains an example REST service that connects FOX API to RabbitMQ using micropipe-proxy.
REST service is written using Node.js

## Requirements

To run the demo you need to have Docker and Docker-Compose installed.

## Running

1. Clone this repo: `git clone git@github.com:AKSW/micropipe-proxy.git`
2. Enter the example folder: `cd ./micropipe-proxy/example/fox`
3. Build and pull required docker images: `docker-compose build`
4. Start the services: `docker-compose up`

## Usage

Once service is up, RabbitMQ instance will be accessible locally on default RabbitMQ port (5672).
To request processing the data by FOX, you need to send your data object that should have field `text` to route `fox-v1`, e.g.:
```json
{
    "route": "fox-v1",
    "data": {
        "text": "The philosopher and mathematician Leibniz was born in Leipzig in 1646 and attended the University of Leipzig from 1661-1666. The current chancellor of Germany, Angela Merkel, also attended this university."
    }
}
```
The micropipe-proxy will reply either to `replyTo` (when specified) or to next part of topic (e.g. if topic if `fox-v1.other-service`, the reply will go to `other-service`).

For testing purposes, you can send input data directly to micropipe-proxy on port 8080.
For example usage see `make test` command.
