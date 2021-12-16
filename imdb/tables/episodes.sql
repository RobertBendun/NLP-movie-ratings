CREATE TABLE "Episodes" (
	"Episode_ID"	TEXT NOT NULL,
	"Show_ID"	TEXT NOT NULL,
	"Season"	INT,
	"Episode" INT,
	PRIMARY KEY("Episode_ID", "Show_ID"),
	FOREIGN KEY(Episode_ID) REFERENCES Basics(ID),
	FOREIGN KEY(Show_ID) REFERENCES Basics(ID)
);
