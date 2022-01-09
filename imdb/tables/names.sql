CREATE TABLE "Names" (
	"ID"	TEXT NOT NULL UNIQUE,
	"Primary_Name"	TEXT NOT NULL,
	"Birth_Year"	INTEGER,
	"Death_Year"	INTEGER,
	"Profession"	TEXT,
	"Known_For"	  TEXT,
	PRIMARY KEY("ID")
);
