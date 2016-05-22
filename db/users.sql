CREATE TABLE users (
	uid       serial PRIMARY KEY,
	user_name varchar(255),
	full_name varchar(255),
	email     varchar(255) NOT NULL,
	password  varchar(255),
	created   timestamp DEFAULT current_timestamp,
	enabled	  boolean DEFAULT TRUE
);

CREATE UNIQUE INDEX email_idx ON users (email);
