default: run

run:
	go run *.go

test:
	echo '{"route": "exynize.test", "data": {"test": "json"}}' | http POST localhost:8080

build:
	go build

rabbit:
	docker run -d -p 5672:5672 -p 8081:15672 --name exynize-rabbit rabbitmq:management

stop-rabbit:
	docker stop exynize-rabbit
