all: imdb imdb-data db-tool/db-tool

db-tool/db-tool: db-tool/*.go
	cd db-tool; go build

# Database creation
imdb: imdb-data db-tool/db-tool imdb-tables/*.sql
	rm -f imdb
	sqlite3 imdb <(cat imdb-tables/*.sql)
	db-tool/db-tool -db imdb -table Ratings -tsv imdb-datasets/title.ratings.tsv
	db-tool/db-tool -db imdb -table Basics -tsv imdb-datasets/title.basics.tsv

# IMDB DATA
datasets=\
    imdb-datasets/name.basics.tsv\
    imdb-datasets/title.akas.tsv\
    imdb-datasets/title.basics.tsv\
    imdb-datasets/title.crew.tsv\
    imdb-datasets/title.episode.tsv\
    imdb-datasets/title.principals.tsv\
    imdb-datasets/title.ratings.tsv

.PHONY: imdb-data
imdb-data: $(datasets)

imdb-datasets/%.tsv: imdb-datasets/%.tsv.gz
	gzip -d $<

imdb-datasets/%.tsv.gz:
	wget 'https://datasets.imdbws.com/$(shell basename $@)' -O $@ -q

.PHONY: clean
clean:
	rm -f db-tool/db-tool

.PHONY: clean-data
clean-data:
	rm -f imdb-datasets/*.tsv imdb-datasets/*.gz

.PHONY: clean-db
clean-db:
	rm -f imdb

.PHONY: clean-all
clean-all: clean clean-data clean-db
