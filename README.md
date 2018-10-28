# Kafka Webhook

Send all records of a kafka topic to http endpoint.

## Run webhook

```bash
go run main.go \
-port=8080 \
-kafka-brokers=kafka:9092 \
-kafka-topic=mytopic \
-kafka-group=mygroup \
-hook-url=http://localhost:1234/hook \
-hook-method=POST \
-secret=DontTellAnybody \
-retry-limit=3 \
-retry-delay=2s \
-v=2
```

## Test setup

Start debug server

```bash
go get github.com/bborbe/debug-server
debug-server -port 1234
```

Produce message

```bash
go get github.com/Shopify/sarama/tools/kafka-console-producer
kafka-console-producer \
-brokers kafka:9092 \
-topic mytopic \
-key myKey \
-value myValue 
```
