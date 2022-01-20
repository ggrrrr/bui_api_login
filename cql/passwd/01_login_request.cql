
CREATE TABLE IF NOT EXISTS request(
	email text,
	attr map<text,text>,
	created_time timestamp,
	enabled boolean,
	PRIMARY KEY (email)
);

