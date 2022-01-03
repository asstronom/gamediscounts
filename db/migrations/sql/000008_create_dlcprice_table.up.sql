CREATE TABLE IF NOT EXISTS dlcprice (
		dlcid INT REFERENCES dlc (id) ON UPDATE CASCADE ON DELETE CASCADE,
		storeid INT REFERENCES store (id) ON UPDATE CASCADE ON DELETE CASCADE,
		CONSTRAINT dlcPriceId PRIMARY KEY (dlcid, storeid),
		storedlcid text UNIQUE,
		initial NUMERIC,
		final NUMERIC, 
		discount INT DEFAULT 0,
		free BOOLEAN DEFAULT FALSE,
		currency VARCHAR (10)
);