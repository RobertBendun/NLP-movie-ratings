ensure_bin=@mkdir -p bin

all: imdb imdb-data \
	raport.pdf \
	bin/metrics bin/psplit bin/db-tool bin/flatten

progs/db-tool/db-tool: progs/db-tool/*.go
	cd $(@D); go build

bin/db-tool: progs/db-tool/db-tool
	$(ensure_bin)
	ln -fs "$(shell pwd)/$<" $@

bin/%: progs/%.cc
	$(ensure_bin)
	$(CXX) -std=c++20 -Wall -Wextra -o $@ $< -lfmt

# Database creation
db: imdb-data bin/db-tool imdb-tables/*.sql
	rm -f $@
	sqlite3 $@ <(cat imdb-tables/*.sql)
	bin/db-tool -db $@ -table Ratings -tsv imdb/data/title.ratings.tsv
	bin/db-tool -db $@ -table Basics -tsv imdb/data/title.basics.tsv

# Raport

%.pdf: %.tex
	pdflatex $<

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
	rm -rf bin
	rm -f progs/db-tool/db-tool
	rm -f *.{log,aux,gdb_latexmk,fls,log,pdf}

.PHONY: clean-data
clean-data:
	rm -f imdb/data/*.tsv imdb/data/*.gz

.PHONY: clean-db
clean-db:
	rm -f db

.PHONY: clean-all
clean-all: clean clean-data clean-db
