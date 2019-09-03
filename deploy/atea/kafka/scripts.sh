## group info
kafka-consumer-groups --bootstrap-server kafka-svc.atea.svc.cluster.local:9092 --describe --group Transcriber.Service.Group

## consume all messages
kafka-console-consumer --bootstrap-server kafka-svc.atea.svc.cluster.local:9092 --group Transcriber.Service.Group --topic NewAudioAvailable
