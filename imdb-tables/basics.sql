CREATE TABLE "Basics" (
	"ID"	TEXT NOT NULL UNIQUE,
	"Type"	TEXT NOT NULL,
	"Primary_Title"	TEXT NOT NULL,
	"Original_Title"	TEXT NOT NULL,
	"IsAdult"	INTEGER NOT NULL,
	"Start_Year"	INTEGER,
	"End_Year"	INTEGER,
	"Runtime"	INTEGER,
	"Genres"	TEXT,
	PRIMARY KEY("ID")
);
