CREATE TABLE IF NOT EXISTS gameprice (
		gameid INT REFERENCES game (id) ON UPDATE CASCADE ON DELETE CASCADE,
		storeid INT REFERENCES store (id) ON UPDATE CASCADE ON DELETE CASCADE,
		CONSTRAINT gamePriceId PRIMARY KEY (gameid, storeid),
		storegameid text UNIQUE,
		initial NUMERIC,
		final NUMERIC, 
		discount INT DEFAULT 0,
		free BOOLEAN DEFAULT FALSE,
		currency VARCHAR (10)
);