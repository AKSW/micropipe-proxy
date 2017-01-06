run:
	./micropipe-proxy-mac

rabbit:
	docker run -d -p 5672:5672 -p 8081:15672 --name geiser-rabbit rabbitmq:management

stop-rabbit:
	docker stop geiser-rabbit

fox:
	docker run -d -p 4444:4444 --name geiser-fox bde2020/fox

stop-fox:
	docker stop geiser-fox

# echo '{"route": "fox-v1", "data": {"text": "The philosopher and mathematician Leibniz was born in Leipzig in 1646 and attended the University of Leipzig from 1661-1666. The current chancellor of Germany, Angela Merkel, also attended this university."}}' | http POST localhost:8080

test:
	http POST localhost:8080 < input.json

test-local:
	echo '{"body": {"text": "The philosopher and mathematician Leibniz was born in Leipzig in 1646 and attended the University of Leipzig from 1661-1666. The current chancellor of Germany, Angela Merkel, also attended this university."}}' | http POST localhost:3000
