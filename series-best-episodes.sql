SELECT 
	Episode_ID AS "ID Odcinka",
	show.Original_Title AS "Serial", 
	episode.Original_Title AS "Odcinek", 
	Rating AS "Ocena", 
	Votes_Count AS "Liczba głosów", 
	Season AS "NR sezonu",
	Episode AS "NR odcinka",
	episode.Start_Year AS "Rok Premiery"
FROM Episodes
	INNER JOIN Basics show ON Show_ID = show.ID
	INNER JOIN Basics episode ON Episode_ID = episode.ID
	INNER JOIN Ratings ON Episodes.Episode_ID = Ratings.Title
WHERE show.Original_Title LIKE "Winx Club"
ORDER BY Rating DESC, Votes_Count DESC