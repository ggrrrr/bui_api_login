
CREATE TABLE IF NOT EXISTS user_passwd(
	email text,
	namespaces set<text>,
	attr map<text,text>,
	created_time timestamp,
	enabled boolean,
	passwd text,
	PRIMARY KEY (email)
);

