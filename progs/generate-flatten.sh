#!/usr/bin/env bash

for file in queries/*.sql; do
	Base="$(basename $file)"
	Base="${base%%.*}.txt"
	Flat="flat-$Base"

	bin/flatten "$Base" 0.1 > "$Flat"
	Base="$Flat"
	Data="$Base.txt"
	Train="$Base.train.txt"
	Test="$Base.test.txt"
	Predictions="$Base.predictions.txt"
	Model="$Base.vw"

	echo "!!! Generating $Base.txt from database"
	bin/db-tool vw -db db -query "$1" -out "$Data"

	echo "!!! Splitting into train and test set"
	bin/psplit 3 "$Train" 1 "$Test" < "$Data"

	echo "!!! Generating model"
	vw -d "$Train" -f "$Model"

	echo "!!! Evaluating"
	vw -d "$Test" -i "$Model" -p "$Predictions"

	echo "!!! Generating plots"
	progs/plot-error.py "$Test" "$Predictions"
done
