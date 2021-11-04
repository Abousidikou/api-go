#!/bin/bash

echo "Data Retrieving..."

## Getting today path
year=`date -d 'yesterday' +%Y`
month=`date -d 'yesterday' +%m`
day=`date -d 'yesterday' +%d`
script_path="/home/emery/ndt/script_ndt_auto_report/"
path="/home/emery/ndt/ndt-server/datadir/ndt7"
log_file="/home/emery/ndt/script_ndt_auto_report/Daily_Report.log"
wdr="/tmp/ndt14QRT589YUSD695LPOKIJSDDFytr/"
header=""

cd $script_path

if [[ $@ == 5 ]]; then
	path=$5
fi


if [[ $# == 5 ]]; then
	path=$path"/"$2"/"$3"/"$4"/"
	year=$2
	month=$3
	day=$4
elif [[ $# == 4 ]]; then
	path=$path"/"$2"/"$3"/"$4"/"
	year=$2
	month=$3
	day=$4
elif [[ $# == 1 ]]; then
	path=$path"/"$year"/"$month"/"$day"/"
else
	echo "Your aurgument is not correct... "
	exit
fi


### Header
header="$year-$month-$day Report..."
if [[ $# == 2 ]]; then
	header="$year-$month-$2 Report..."	
fi
echo  $header
echo  $header >> $log_file
## Check path or return

if [ ! -d $path ]
then
	echo $path
	echo "No test was made "
	echo -e "No test was made  \n\n\n" >> $log_file
	exit
fi

echo "Directoty: "$path
echo "Directoty: "$path >> $log_file

## Creating tmp file tmp/ndt

if [ ! -d $wdr ]
then
	mkdir $wdr
	echo  "Creation of directory $wdr" >> $log_file
else
	rm -r $wdr
	mkdir $wdr
	echo "Creation of directory $wdr " >> $log_file
fi


## Copy file into tmp/ndt directory

for file in `ls $path`;
do
	extension=`echo $file | rev | cut -d'.' -f1 | rev`
	if [[ $file != "zip" && $file != "unzip" ]]; then
		echo  "Copying  $path$file " | tee --append $log_file
		cp $path$file  $wdr
	fi
done


## Gzipping decompressing file
for file in `ls $wdr`; do
	extension=`echo $file | rev | cut -d'.' -f1 | rev`
	if [[ $extension == "gz" ]]; then
		echo  "Dezipping  $wdr$file " | tee --append $log_file
		gzip -d $wdr$file
	fi
done


## Start getting Data

if [[ $1 == "setup" ]]; then
	echo  "Setting Database..." | tee --append $log_file
	python3 process_pro.py "setup"
	echo -e "Database Set  \n\n\n" | tee --append $log_file
fi

if [[ $1 == "report" ]]; then
	echo  "Daily push into database..." | tee --append $log_file
	for file in `ls $wdr`;
	do
		if [[ $file != "zip" ]]; then
			echo  "Processing  $file " | tee --append $log_file
			python3 process_pro.py $wdr$file
			echo "Processing of  $file finished" | tee --append $log_file
		fi
	done
	# dayly Statistic
	echo "Daily stat... " | tee --append $log_file
	if [[ $# == 2 ]]; then
		python3 process_pro.py "dayStat" $year $month $2
	else
		python3 process_pro.py "dayStat" $year $month $day
	fi 
	echo "Daily stat finished" | tee --append $log_file
fi




## sleep 5
echo "Deletion of directory  $wdr" | tee --append $log_file
rm -r $wdr
echo -e  "End of today's report  $wdr \n\n\n" | tee --append $log_file