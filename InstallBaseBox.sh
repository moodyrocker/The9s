#!/bin/bash


sudo apt-get update && sudo apt-get install jq golang-go influxdb docker.io influxdb-client sqlite3 pwgen

sudo echo 500000 >/proc/sys/fs/file-max



#58  export GOPATH=/home/ubuntu/The9s/go
#65  go get -u -v "github.com/julienschmidt/httprouter"
#*/10 * * * * /home/ubuntu/The9s/ReconScript.sh > /homr/ubuntu/recon.log 2>&1
#back to home
#cd
#git clone https://github.com/mudassar-mian/The9s.git
#mkdir -p /root/The9s/go/src/statsapp/templates/tstats/
#cp The9s/telegraf/telegraf.conf /etc/telegraf/
