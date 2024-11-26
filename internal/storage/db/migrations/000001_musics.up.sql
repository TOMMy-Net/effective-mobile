CREATE TABLE IF NOT EXISTS "musics" (
	"id" serial NOT NULL UNIQUE,
	"song" VARCHAR(255) NOT NULL,
	"music_group"  VARCHAR(255) NOT NULL,
	"text" TEXT NOT NULL DEFAULT '',
	"link" VARCHAR(255) NOT NULL DEFAULT '',
	"release_date" DATE NOT NULL,
	PRIMARY KEY("id")
);

