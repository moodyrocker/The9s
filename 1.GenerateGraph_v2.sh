#!/bin/bash

#Fixed file paths
rm -fr $(pwd)/go/templates/stats/*
Gendate=`date +"%H:%M %d/%m/%Y"`
IFS=$'\n'

# data format "website,groupID" sting test blash
for account_id in `cat 100TopSites.file | awk '{print $2}' | sort -u `
	do
		counter=1
			echo "Now getting data for account with ID > " $account_id

			for url_id in `cat 100TopSites.file|  awk '$2 == '$account_id' { print $0 }' | awk '{print $3}'`
			do

				echo "Account ID...... " $account_id
				echo "Now getting data for URL with ID.... " $url_id
				for Timecheck in 6h 1d 7d 30d
					do
					# This bit will sort out the group by in the Main graph data for each time span
						case $Timecheck in
							6h)
							GroupbyDate=15m
							TimeFrameFormat='6 hours';;
							1d)
							GroupbyDate=30m
							TimeFrameFormat='24 hours';;
							7d)
							GroupbyDate=180m
							TimeFrameFormat='7 days';;
							30d)
							GroupbyDate=1d
							TimeFrameFormat='30 days';;
						esac
							DataFile=$(pwd)'/go/templates/stats/'$account_id'_'$Timecheck'.html'
							echo "Account ID...... " $account_id
							echo "Now getting data for URL with ID...." $url_id
							#the number of checks and number good responce = then give ths percentage in the time period per graph
							CheckCount=`curl -s -G 'http://localhost:8086/query?pretty=true' --data-urlencode "db=BBstats" --data-urlencode "q=SELECT count(*) FROM http_response_times where url_id = '$url_id' AND u_id = '$account_id' AND time > now() - $Timecheck" | jq -c -r .results[0].series[0].values[0][1]`
							echo "Now getting data for" $Timecheck
							echo "The of check found for this period" $CheckCount

							# sleep 3

							CountGoodResponces=`curl -s -G 'http://localhost:8086/query?pretty=true' --data-urlencode "db=BBstats" --data-urlencode "q=SELECT count(*) FROM http_response_times where url_id = '$url_id' AND u_id = '$account_id' AND (ServerStatusCode =~ /30/ OR ServerStatusCode =~ /20/) AND time > now() - $Timecheck" | jq -c -r  .results[0].series[0].values[0][1]`

							UptimeStats=`awk "BEGIN {print ($CountGoodResponces/$CheckCount)*100}" | awk '{printf("%.2f\n", $1)}'`
							echo "uptime >" $UptimeStats
							echo '{{define "UptimeStats_'$counter'"}}'$UptimeStats'{{end}}' >>$DataFile

							#AVG Responce time for a given period, 1 7 30 days
							echo '{{define "AVGResponseTime_'$counter'"}}'>>$DataFile
							curl -s -G 'http://localhost:8086/query?pretty=true' --data-urlencode "epoch=ms" --data-urlencode "db=BBstats" --data-urlencode "q=SELECT MEAN(value) FROM http_response_times where url_id = '$url_id' AND u_id = '$account_id' AND time > now() - $Timecheck limit 1" | jq -c -r  .results[].series[].values[][1] | awk '{printf "%.2f", $1}'>>$DataFile
							echo '{{end}}'>>$DataFile

							HTTPResonseTimeData=`curl -s -G 'http://localhost:8086/query?pretty=true' --data-urlencode "epoch=ms" --data-urlencode "db=BBstats" --data-urlencode "q=SELECT MEAN(value) FROM http_response_times where url_id = '$url_id' AND u_id = '$account_id' AND time > now() - $Timecheck group by time($GroupbyDate)" | jq -c -r .results[0].series[0].values`

							get_short_URL=`sqlite3 $(pwd)/the9s "select short_url from url where url_id=$url_id and u_id =$account_id;"`
							get_URL=`sqlite3 $(pwd)/the9s "select url from url where url_id=$url_id and u_id =$account_id;"`

							echo '{{define "GraphTitle_'$counter'"}}'$get_short_URL'{{end}}' >>$DataFile
							echo '{{define "GraphTitleLink_'$counter'"}}'$get_URL'{{end}}' >>$DataFile

							echo '{{define "CheckCount_'$counter'"}}'$CheckCount'{{end}}'>>$DataFile
							echo '{{define "GraphColumnchart_'$counter'"}}' >>$DataFile
							echo $HTTPResonseTimeData >>$DataFile
							echo '{{end}}'>>$DataFile
							echo '{{define "TimeCheck_'$counter'"}}'$Timecheck'{{end}}' >>$DataFile
							echo '{{define "TimeFrameFormat_'$counter'"}}'$TimeFrameFormat'{{end}}' >>$DataFile

						done
					counter=$((counter+1))
					echo "******************"
	done
done
