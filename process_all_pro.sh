#!/bin/bash

## Getting today path
path="/home/emery/ndt/ndt-server/datadir/ndt7"
for year in `ls $path`; do
	for  month in `ls $path"/"$year`; do
		for day in `ls $path$"/"$year"/"$month`; do
			./process_pro.sh "report" $year $month $day  $path
		done
	done
done