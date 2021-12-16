CREAtE TABLE "Localization" (
	"Title_ID" TEXT NOT NULL,
	"Ordering" INT NOT NULL,
	"Title" TEXT NOT NULL,
	"Region" TEXT,
	"Language" TEXT,
	"Types" TEXT,
	"Attributes" TEXT,
	"IsOriginal" INT,

	PRIMARY KEY("Title_ID", "Ordering"),
	FOREIGN KEY("Title_ID") REFERENCES Basics(ID)
);
