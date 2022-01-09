SELECT DISTINCT
	Rating AS 'float',
	Original_Title AS 'str'
FROM Basics
INNER JOIN Ratings ON Basics.ID == Ratings.Title;
