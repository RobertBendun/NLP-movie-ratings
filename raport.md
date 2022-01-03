---
title:
---
# Tytuły jako wyłączny wyznacznik oceny filmu

Dane wyznaczone przy pomocy następującej kwerendy:
```sql
SELECT
	CAST(ROUND(Rating) AS INT) AS 'nat',
	Original_Title AS 'str'
FROM Basics
INNER JOIN Ratings ON Basics.ID == Ratings.Title;
```

A następnie przygotowano dane w następujący sposób:

```console
$ ./db-tool/db-tool vw -db db -query titles.sql -out titles.txt
$ ./pslit 3 titles.train.txt 1 titles.test.1.txt <titles.txt
$ cut -d' ' -f1 --complement titles.test.1.txt titles.test.txt
```

Co wygenerowało dwa zbiory: uczący oraz testowy.

```console
$ wc -l titles.{test,train}.txt
  237930 titles.test.txt
  716976 titles.train.txt
  954906 total
```

Pozwala to wygenerować model:

```console
$ vw -d titles.train.txt -f titles.vw
```
