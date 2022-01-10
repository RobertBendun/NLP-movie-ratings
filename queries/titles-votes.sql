SELECT
	Rating AS 'id',
	Original_Title AS 'bag',
	Votes_Count AS 'prefix("__votes:")'
FROM Basics
	INNER JOIN Ratings ON Basics.ID == Ratings.Title;
