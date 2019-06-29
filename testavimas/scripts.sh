## info
https://app.swaggerhub.com/apis/aireno/Transkipcija/1.1.1

## system
cat /etc/*-release

## proxy
kubectl -n aft port-forward svc/upload 8001:8000
kubectl -n aft port-forward svc/status 8002:8000
kubectl -n aft port-forward svc/result 8003:8000

## live 
curl -i http://localhost:8001/live?full=1
curl -i http://localhost:8002/live?full=1
curl -i http://localhost:8003/live?full=1

## data
curl -k https://prn509.vdu.lt:7080/testdata/testdata.zip -o testdata.zip
unzip testdata.zip

## 1.1
curl -X POST -i http://localhost:8001/upload -H 'Accept: application/json' -H 'content-type: multipart/form-data' -F file=@testdata/wav/t_aft_11.wav

## 1.2
curl -X POST -i http://localhost:8001/upload -H 'Accept: application/json' -H 'content-type: multipart/form-data'

## 2.1
curl -X POST -i http://localhost:8001/upload -H 'Accept: application/json' -H 'content-type: multipart/form-data' -F file=@testdata/wav/t_aft_11.wav
curl -i http://localhost:8002/status/

## 2.2
curl -X POST -i http://localhost:8001/upload -H 'Accept: application/json' -H 'content-type: multipart/form-data' -F file=@testdata/wav/t_aft_22.wav
curl -i http://localhost:8002/status/

## 3.1
curl http://localhost:8003/result/{ID}/result.txt -o result.txt
vimdiff testdata/txt/t_aft_11.txt result.txt 


## 3.2
curl http://localhost:8003/audio/{ID} -o out.wav
play out.wav
play testdata/wav/t_aft_11.wav

## 4.1
ls -1 testdata/wav/t_aft_4*.wav | xargs -n1 -P10 ./bin/send.sh > fl
cat fl | xargs -n2 -P20 ./bin/status.sh | sort
cat fl | xargs -n2 -P20 ./bin/result.sh
./bin/test_wer.sh 
vimdiff testdata/txt/ref.txt recognized.txt 

