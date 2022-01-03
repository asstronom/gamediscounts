CREATE TABLE IF NOT EXISTS dlc (
		id SERIAL PRIMARY KEY,
		gameid INT REFERENCES game (id) ON UPDATE CASCADE ON DELETE CASCADE,
		name text UNIQUE,
        shortdescription TEXT,
		description TEXT,
		headerImage TEXT DEFAULT 'https:\/\/cdn.akamai.steamstatic.com\/steam\/apps\/292030\/header_russian.jpg?t=1621939214'
);