#!/bin/bash

## this builds the top part of the stats page
#
Gendate=`date +"%H:%M %d/%m/%Y"`
# data format "website,groupID" sting test blash
for account_id in `cat 100TopSites.file | awk '{print $2}' | sort -u `
do
# Grouping websites together then generating the data based on groups
  counter=1
    for url_nf in `grep -w $account_id 100TopSites.file`
      do

        url=`echo $url_nf | awk  '{print $1}'`

        MAINDATAFILE=$(pwd)'/go/templates/stats/'$account_id'_main.html'
        > $MAINDATAFILE

        account_name=`sqlite3 /home/vagrant/The9s/the9s "select account_name from user where u_id=$account_id;"`
        echo '{{define "SummeryTitle"}}'$account_name'{{end}}' >> $MAINDATAFILE
	       #### CurrentHTTPResponseCode ####
	       CurrentHTTPResponseCode=`curl -s -G 'http://localhost:8086/query?pretty=true' --data-urlencode "epoch=ms" --data-urlencode "db=BBstats" --data-urlencode "q=SELECT value,ServerStatusCode FROM http_response_times where url_id = '$url' order by desc limit 1" | jq -c -r  .results[0].series[0].values[0][2]`

         if [[ $CurrentHTTPResponseCode =~ ^2 ]] || [[ $CurrentHTTPResponseCode =~ ^3 ]]
            then
              chartbgcolor=""
              mainstatus="<span class='uk-text-large uk-text-success'>All systems are good!</span>"
            else
              chartbgcolor="alerterror"
              mainstatus="<span class='alerterror uk-text-danger uk-text-large'>Some systems have degraded</span>"
            fi

	           echo '{{define "MainStatus"}}'$mainstatus'{{end}}' >>$MAINDATAFILE
             echo '{{define "Summeryaccount_id"}}'$account_id'{{end}}'>>$MAINDATAFILE
	           echo '{{define "LastTimeCheck"}}'$Gendate'{{end}}' >>$MAINDATAFILE
            counter=$((counter+1))

            echo '{{define "info"}}'  >>$MAINDATAFILE
            echo '
              <li>
                <div uk-grid style="font-size:14px;">
                <div class="uk-width-1-6" > <span uk-icon="icon: warning; ratio: 1.2" style="color:red"></span></div>
                <div class="uk-width-1-6" ><b>16:30 14/01/2019</b></div>
                <div class="uk-width-1-6" ><b>www.ibm.com</b></div>
                <div class="uk-width-expand"><p>503 error have been seen at this time. For more information click on the link below</p></div>
                </div>
              </li>
              <li>
                <div uk-grid style="font-size:14px;">
                <div class="uk-width-1-6" > <span uk-icon="icon: warning; ratio: 1.2" style="color:red"></span></div>
                <div class="uk-width-1-6" ><b>16:30 14/01/2019</b></div>
                <div class="uk-width-1-6" ><b>www.ibm.com</b></div>
                <div class="uk-width-expand"><p>503 error have been seen at this time. For more information click on the link below</p></div>
                </div>
              </li>'>>$MAINDATAFILE
            echo '{{end}}' >>$MAINDATAFILE
            ###### end ######
      done
  done
