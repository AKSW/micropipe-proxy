default: run

run:
	go run *.go

test:
	echo '{"route": "exynize.test", "data": {"test": "json"}}' | http POST localhost:8080

build:
	go build

build-linux:
	GOOS=linux GOARCH=amd64 go build -o proxy-linux

rabbit:
	docker run -d -p 5672:5672 -p 8081:15672 --name exynize-rabbit rabbitmq:management

stop-rabbit:
	docker stop exynize-rabbit
