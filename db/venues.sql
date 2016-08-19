CREATE TABLE venues (
	vid       serial PRIMARY KEY,
	name      varchar(255) NOT NULL,
	address   varchar(255),
	coords    varchar(255),
);

CREATE UNIQUE INDEX name_idx ON venues (name);
