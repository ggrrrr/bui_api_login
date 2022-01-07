
CREATE TABLE IF NOT EXISTS role (
	namespace TEXT,
	name TEXT,
	services frozen<set<text>>,
	resources set<text>,
	verbs set<text>,
	created_time timestamp,
	PRIMARY KEY ((namespace,name), services)
);