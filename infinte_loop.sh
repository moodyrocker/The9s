#!/bin/bash
while :
do
	./Get_URL_Stats.sh && ./1.GenerateGraph_v2.sh && ./2.Generate_js_v2.sh && ./3.GenerateMain_v2.sh
	sleep 60
done
