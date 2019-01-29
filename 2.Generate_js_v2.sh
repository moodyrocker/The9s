#!/bin/bash

ChartColours=(00000 ECD078 D95B43 C02942 542437 53777A 3C3243 6B4880)
for account_id in `cat 100TopSites.file | awk '{print $2}' | sort -u `
  do
    echo "Account id" $account_id
	  echo '{{define "jsforcharts"}}' > $(pwd)'/go/templates/stats/chartjs'$account_id'.inc.go.html'
		echo '{{define "subsections"}}' > $(pwd)'/go/templates/stats/subsec'$account_id'.inc.go.html'

      count=1
      for wbf in `cat 100TopSites.file|  awk '{print $2}' | grep -w $account_id`
        do
          echo $wbf
				    cat blankjs.txt | sed -e 's/XXXBGColourXXX/'$chartbgcolor'/g' -e 's/XXXcontainerXXX/container'$count'/g' -e 's/XXXcolourXXX/'${ChartColours[$count]}'/g' -e 's/XXXdataXXX/{{template "GraphColumnchart_'$count'"}}/g' >>$(pwd)'/go/templates/stats/chartjs'$account_id'.inc.go.html'
			      cat blanksection.txt | sed -e 's/XXXstatusXXX/'$chartbgcolor'/g' -e 's/XXXcolourXXX/'${ChartColours[$count]}'/g' -e 's/AXAA/'$count'/g' >>$(pwd)'/go/templates/stats/subsec'$account_id'.inc.go.html'
          	count=$((count+1))
           echo "count" $count
        done

		        echo '{{end}}' >> $(pwd)/go/templates/stats/chartjs${account_id}.inc.go.html
		        echo '{{end}}' >> $(pwd)/go/templates/stats/subsec${account_id}.inc.go.html
  done
