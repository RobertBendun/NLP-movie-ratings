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

```
:!g++ metrics.cc -o metrics -O3 -lfmt && ./metrics titles.test.txt predictions.txt
Accuracy: 30.21% (equal: 71 887, total: 237 930)
Precision:
  1: 0.00% (equal: 0, total: 551)
  2: 0.00% (equal: 0, total: 1 295)
  3: 0.00% (equal: 0, total: 3 384)
  4: 14.29% (equal: 3, total: 8 506)
  5: 24.07% (equal: 84, total: 19 810)
  6: 24.16% (equal: 2 296, total: 42 060)
  7: 29.91% (equal: 56 972, total: 70 107)
  8: 33.68% (equal: 12 306, total: 66 622)
  9: 20.71% (equal: 205, total: 22 394)
  10: 43.75% (equal: 21, total: 3 201)
  Average: 19.06%
Recall:
  1: 0.00% (equal: 0, total: 551)
  2: 0.00% (equal: 0, total: 1 295)
  3: 0.00% (equal: 0, total: 3 384)
  4: 0.04% (equal: 3, total: 8 506)
  5: 0.42% (equal: 84, total: 19 810)
  6: 5.46% (equal: 2 296, total: 42 060)
  7: 81.26% (equal: 56 972, total: 70 107)
  8: 18.47% (equal: 12 306, total: 66 622)
  9: 0.92% (equal: 205, total: 22 394)
  10: 0.66% (equal: 21, total: 3 201)
  Average: 10.72%
```
