
CREATE TABLE IF NOT EXISTS namespaces (
	namespace TEXT,
	domain TEXT,
	created_time timestamp,
	PRIMARY KEY ((namespace), domain)
)

 