CREATE TABLE users (
	uid       serial,
	user_name varchar(255),
	full_name varchar(255),
	email     varchar(255),
	password  varchar(255),
	created   timestamp,
	enabled	  boolean
);
