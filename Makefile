all: imdb

.PHONY: clean
clean:

.PHONY: clean-all
clean-all: clean
	rm -f imdb-datasets/*.tsv imdb-datasets/*.gz

# IMDB DATA
datasets=\
    imdb-datasets/name.basics.tsv\
    imdb-datasets/title.akas.tsv\
    imdb-datasets/title.basics.tsv\
    imdb-datasets/title.crew.tsv\
    imdb-datasets/title.episode.tsv\
    imdb-datasets/title.principals.tsv\
    imdb-datasets/title.ratings.tsv

.PHONY: imdb
imdb: $(datasets)

imdb-datasets/%.tsv: imdb-datasets/%.tsv.gz
	gzip -d $<

imdb-datasets/%.tsv.gz:
	wget 'https://datasets.imdbws.com/$(shell basename $@)' -O $@ -q
