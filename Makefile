all: imdb imdb-data db-tool/db-tool raport.pdf

db-tool/db-tool: db-tool/*.go
	cd db-tool; go build

# Database creation
db: imdb-data db-tool/db-tool imdb-tables/*.sql
	rm -f $@
	sqlite3 $@ <(cat imdb-tables/*.sql)
	db-tool/db-tool -db $@ -table Ratings -tsv imdb/data/title.ratings.tsv
	db-tool/db-tool -db $@ -table Basics -tsv imdb/data/title.basics.tsv

# Raport

%.pdf: %.md
	pandoc -Tpdf -o $@ $<

# IMDB DATA
datasets=\
    imdb/data/name.basics.tsv\
    imdb/data/title.akas.tsv\
    imdb/data/title.basics.tsv\
    imdb/data/title.crew.tsv\
    imdb/data/title.episode.tsv\
    imdb/data/title.principals.tsv\
    imdb/data/title.ratings.tsv

.PHONY: imdb-data
imdb-data: $(datasets)

imdb/data/%.tsv: imdb/data/%.tsv.gz
	gzip -d $<

imdb/data/%.tsv.gz:
	wget 'https://datasets.imdbws.com/$(shell basename $@)' -O $@ -q

.PHONY: clean
clean:
	rm -f db-tool/db-tool

.PHONY: clean-data
clean-data:
	rm -f imdb/data/*.tsv imdb/data/*.gz

.PHONY: clean-db
clean-db:
	rm -f db

.PHONY: clean-all
clean-all: clean clean-data clean-db
