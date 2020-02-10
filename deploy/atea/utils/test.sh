#/bin/bash
######################
proxyTask=$1
testTask=$2
######################
echo Creating proxy
make $proxyTask &
pid=$!
echo Created proxy. Task: $pid
######################
echo Wait
sleep 1
######################
echo Testing
make $testTask
echo Tested
######################
echo Stop proxy
ps --ppid $pid | awk '{print $1}' | tail -n +2 | xargs kill -9
echo Stopped
######################