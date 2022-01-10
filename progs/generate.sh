#!/usr/bin/env bash

set -e

Base="$(basename "$1")"
Base="${Base%%.*}"

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
