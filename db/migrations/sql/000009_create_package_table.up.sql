CREATE TABLE IF NOT EXISTS package (
		id SERIAL PRIMARY KEY,
		name TEXT UNIQUE,
		smalllogo TEXT,
		pageimage TEXT
);