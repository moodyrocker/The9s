#!/bin/bash

echo "running the test should take 30s...."

IFS=$'\n'
for i in `cat $HOME/The9s/100TopSites.file`
do
	url=`echo $i | awk '{print $1}'`
	u_id=`echo $i | awk '{print $2}'`
	url_id=`echo $i | awk '{print $3}'`
	dirrandom=`pwgen 30 1`

	time sudo docker run --shm-size=1g --rm -v "$(pwd)":/browsertime sitespeedio/browsertime --timeouts.pageCompleteCheck 100000 -q --skipHar --video false --resultDir _results/$dirrandom --iterations  1 $url;

	pageLoadTime=`cat $(pwd)/_results/$dirrandom/browsertime.json | jq -r .statistics.timings.pageTimings.pageLoadTime.mean`
	pageLoadTime_f=`echo $pageLoadTime / 1000 | bc -l | awk '{printf "%.3f", $1}'`
  serverResponseTime=`cat $(pwd)/_results/$dirrandom/browsertime.json | jq -r .statistics.timings.pageTimings.serverResponseTime.mean`
  serverResponseTime_f=`echo $serverResponseTime / 1000 | bc -l | awk '{printf "%.2f", $1}'`
	serverConnectionTime=`cat $(pwd)/_results/$dirrandom/browsertime.json | jq -r .statistics.timings.pageTimings.serverConnectionTime.mean`
	serverConnectionTime_f=`echo $serverConnectionTime / 1000 | bc -l | awk '{printf "%.3f", $1}'`
	fullyLoaded=`cat $(pwd)/_results/$dirrandom/browsertime.json | jq -r .statistics.timings.fullyLoaded.mean`
	fullyLoaded_f=`echo $fullyLoaded / 1000 | bc -l | awk '{printf "%.3f", $1}'`
	RumSpeedIndex=`cat $(pwd)/_results/$dirrandom/browsertime.json | jq -r -c .statistics.timings.rumSpeedIndex.mean`
	ServerStatusCode=`curl -L -s -o /dev/null -w "%{http_code}" $url`

echo "*****************"
echo url: $url
echo url id: $url_id
echo u_id: $u_id
echo pageLoadTime: $pageLoadTime_f
echo serverResponseTime: $serverResponseTime_f
echo serverConnectionTime: $serverConnectionTime_f
echo fullyLoaded: $fullyLoaded_f
echo RumSpeedIndex: $RumSpeedIndex
echo ServerStatusCode: $ServerStatusCode
echo "*****************"

curl -i -XPOST 'http://localhost:8086/write?db=BBstats' --data-binary "http_response_times,url_id=$url_id,u_id=$u_id,ServerStatusCode=$ServerStatusCode,RumSpeedIndex=$RumSpeedIndex,fullyLoaded=$fullyLoaded_f,serverConnectionTime=$serverConnectionTime_f,serverResponseTime=$serverResponseTime_f value=$pageLoadTime_f"

sudo rm -fr $(pwd)/_results/$dirrandom
done
