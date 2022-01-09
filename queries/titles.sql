SELECT DISTINCT
	Rating AS 'id',
	Original_Title AS 'bag'
FROM Basics
INNER JOIN Ratings ON Basics.ID == Ratings.Title;
