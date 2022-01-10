SELECT
	Rating AS 'id',
	Genres AS 'bag',
	Start_Year AS 'id'
FROM
	Basics
	INNER JOIN Ratings ON Basics.ID == Ratings.Title
WHERE
	Genres IS NOT NULL
	AND Start_Year IS NOT NULL;
