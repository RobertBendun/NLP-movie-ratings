SELECT
	Rating AS 'head(id)',
	Person AS 'limitedJoin(" ", id, 3)',
	ID AS 'group',
	Category AS 'limitedCountedJoin(" ", id, 3)'
FROM
	Basics
	INNER JOIN Ratings ON Basics.ID == Ratings.Title
	INNER JOIN Principals ON Basics.ID == Principals.Title
WHERE
	Category IN ("actor", "actress");
