SELECT
	Rating AS 'head(id)',
	Person AS 'join(" ", bag)',
	ID AS 'group'
FROM
	Basics
	INNER JOIN Ratings ON Basics.ID == Ratings.Title
	INNER JOIN Principals ON Basics.ID == Principals.Title;
