SELECT
	ID AS 'str',
	Rating AS 'float'
FROM Basics
INNER JOIN Ratings ON Basics.ID == Ratings.Title
INNER JOIN Principals ON Basics.ID == Principals.Title;
