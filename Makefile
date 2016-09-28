default: run

.PHONY: rabbit

run:
	go run *.go

test:
	echo '{"route": "test-v1", "replyTo": "test-v1", "config": {"test123": {"param": "ok"}}, "data": {"text": "ok"}}' | http POST localhost:8080

test-nodata:
	echo '{"route": "test-v1", "replyTo": "test-v1", "config": {"test123": {"param": "ok"}}}' | http POST localhost:8080

test-noconfig:
	echo '{"route": "test-v1", "replyTo": "test-v1", "data": {"text": "ok"}}' | http POST localhost:8080

test-noinput:
	echo '{"route": "test-v1", "replyTo": "test-v1"}' | http POST localhost:8080

test-sentiments:
	echo '{"route": "sentiments-v1.test-v1", "replyTo": "test-v1", "config": {"sentiments": {"test": "ok"}}, "data": {"text": "I am very awesome text!"}}' | http POST localhost:8080

test-nested:
	echo '{"route": "sentiments-v1.test-v1", "config": {"sentimentConfig": {"one": "ok"}}, "data": {"text": "123", "subobj": {"test": "javascript"}}}' | http POST localhost:8080

test-health:
	http GET localhost:8080/health

build:
	go build

build-linux:
	GOOS=linux GOARCH=amd64 go build -o micropipe-proxy-linux

rabbit:
	docker run -d -p 5672:5672 -p 8081:15672 --name test-rabbit rabbitmq:management

stop-rabbit:
	docker stop test-rabbit
