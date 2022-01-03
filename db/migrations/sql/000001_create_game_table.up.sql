CREATE TABLE IF NOT EXISTS game (
		id SERIAL PRIMARY KEY,
		name TEXT UNIQUE,
		shortdescription TEXT,
		description TEXT,
		headerimage TEXT DEFAULT 'https:\/\/cdn.akamai.steamstatic.com\/steam\/apps\/292030\/header_russian.jpg?t=1621939214'
);