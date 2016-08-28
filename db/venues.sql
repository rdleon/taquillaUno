CREATE TABLE venues (
	vid       serial PRIMARY KEY,
	name      varchar(255) NOT NULL,
	cityId    integer references cities(cityId),
	address   varchar(255),
	coords    varchar(255)
);

CREATE UNIQUE INDEX vname_idx ON venues (name);
