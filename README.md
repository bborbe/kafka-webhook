# Kafka Webhook

Send all records of a kafka topic to http endpoint.

## Run webhook

```bash
go run main.go \
-kafka-brokers=kafka:9092 \
-kafka-topic=mytopic \
-kafka-group=mygroup \
-url=http://www.example.com \
-v=2
```
