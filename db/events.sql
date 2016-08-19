CREATE TABLE events (
	eid       serial PRIMARY KEY,
	name      varchar(255) NOT NULL,
	created	   timestamp DEFAULT current_timestamp,
	start      timestamp with timezone,
	duration   integer, -- Minutes
	published  boolean DEFAULT TRUE
);

CREATE UNIQUE INDEX name_idx ON events (name);
