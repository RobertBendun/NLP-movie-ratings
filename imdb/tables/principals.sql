CREATE TABLE "Principals" (
	"Title"	TEXT NOT NULL,
	"Ordering" INT NOT NULL,
	"Person" TEXT NOT NULL,
	"Category" TEXT NOT NULL,
	"Job" TEXT,
	"Characters" TEXT,
	PRIMARY KEY("Title", "Ordering"),
	FOREIGN KEY("Title") REFERENCES Basics(ID)
);
