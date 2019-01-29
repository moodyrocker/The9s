#!/bin/bash

echo "running the test should take 30s...."

for i in `cat  $HOME/The9s/100TopSites.file`
do
	wb=`echo $i | awk -F',' '{print $1}'`
	wbid=`echo $i | awk -F',' '{print $2}'`

	time sudo docker run --shm-size=1g --rm -v "$(pwd)":/browsertime sitespeedio/browsertime --skipHar --headless --video false --iterations 1 http://www.$wb;

	pageLoadTime=`cat $HOME/The9s/browsertime-results/www.$wb/*/browsertime.json | jq -r .statistics.timings.pageTimings.pageLoadTime.mean`
	pageLoadTime_f=`echo $pageLoadTime / 1000 | bc -l | awk '{printf "%.3f", $1}'`
	serverResponseTime=`cat $HOME/The9s/browsertime-results/www.$wb/*/browsertime.json | jq -r .statistics.timings.pageTimings.serverResponseTime.mean`
	serverResponseTime_f=`echo $serverResponseTime / 1000 | bc -l | awk '{printf "%.3f", $1}'`
	serverConnectionTime=`cat $HOME/The9s/browsertime-results/www.$wb/*/browsertime.json | jq -r .statistics.timings.pageTimings.serverConnectionTime.mean`
	serverConnectionTime_f=`echo $serverConnectionTime / 1000 | bc -l | awk '{printf "%.3f", $1}'`
	fullyLoaded=`cat $HOME/The9s/browsertime-results/www.$wb/*/browsertime.json | jq -r .statistics.timings.fullyLoaded.mean`
	fullyLoaded_f=`echo $fullyLoaded / 1000 | bc -l | awk '{printf "%.3f", $1}'`
	RumSpeedIndex=`cat $HOME/The9s/browsertime-results/www.$wb/*/browsertime.json | jq -r -c .statistics.timings.rumSpeedIndex.mean`
	ServerStatusCode=`curl -L -s -o /dev/null -w "%{http_code}" http://www.$wb`

echo "*****************"
echo $wb
echo $wbid
echo $pageLoadTime_f
echo $serverResponseTime_f
echo $serverConnectionTime_f
echo $fullyLoaded_f
echo $RumSpeedIndex
echo $ServerStatusCode
echo "*****************"


curl -i -XPOST 'http://localhost:8086/write?db=BBstats' --data-binary "http_response_times,wbid=$wbid,website=$wb,ServerStatusCode=$ServerStatusCode,RumSpeedIndex=$RumSpeedIndex,fullyLoaded=$fullyLoaded_f,serverConnectionTime=$serverConnectionTime_f,serverResponseTime=$serverResponseTime_f value=$pageLoadTime_f"

rm -fr $HOME/The9s/browsertime-results/www.$i/*
done
