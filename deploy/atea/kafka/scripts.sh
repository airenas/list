## group info
kafka-consumer-groups --bootstrap-server kafka-svc.atea.svc.cluster.local:9092 --describe --group Transcriber.Service.Group

## consume all messages
kafka-console-consumer --bootstrap-server kafka-svc.atea.svc.cluster.local:9092 --group Transcriber.Service.Group --topic NewAudioAvailable

## show log
 watch "kubectl logs $(kubectl get po | grep -e '^kafka' | tail -n 3| head -n 1 | awk '{print $1}') | tail -n 50"
