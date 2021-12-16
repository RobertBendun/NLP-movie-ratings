SELECT 
	Episode_ID AS "ID Odcinka",
	show.Original_Title AS "Serial",
	Season AS "NR sezonu",
	printf("%.2f", AVG(Rating)) AS "Średnia ocena sezonu"
FROM Episodes
	INNER JOIN Basics show ON Show_ID = show.ID
	INNER JOIN Basics episode ON Episode_ID = episode.ID
	INNER JOIN Ratings ON Episodes.Episode_ID = Ratings.Title
WHERE Season IS NOT NULL AND show.Original_Title LIKE "Winx Club"
GROUP BY Season
ORDER BY AVG(Rating) DESC, Votes_Count DESC