CREATE TABLE IF NOT EXISTS store (
		id SERIAL PRIMARY KEY,
		name text UNIQUE,
		icon text,
		banner text
);