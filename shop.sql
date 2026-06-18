CREATE TABLE products (
	id TEXT PRIMARY KEY,
	name TEXT NOT NULL,
	product_name TEXT NOT NULL,
	thumbnail_url TEXT NOT NULL,
	description TEXT NOT NULL,
	is_ignored BOOLEAN NOT NULL,
	date_created TIMESTAMP NOT NULL,
	date_updated TIMESTAMP NOT NULL,
	files JSONB NOT NULL,
	variant_ids TEXT[] NOT NULL,
	external_id_1 TEXT NOT NULL,
	external_id_2 TEXT NOT NULL,
	extra_data JSONB,
	options JSONB,
	status TEXT NOT NULL
);

CREATE TABLE mockup_tasks (
	id bigserial PRIMARY KEY,
	product_ids TEXT[] NOT NULL,
	source_image INT NOT NULL,
	template JSONB,
	date_created TIMESTAMP NOT NULL,
	date_updated TIMESTAMP NOT NULL,
	status TEXT NOT NULL
);

CREATE TABLE retail_prices (
	product_id TEXT NOT NULL,
	currency TEXT NOT NULL,
	retail_price DECIMAL NOT NULL,
	date_updated TIMESTAMP NOT NULL,
	PRIMARY KEY (product_id, currency)
);

CREATE TABLE images (
	filename TEXT PRIMARY KEY,
	date_created TIMESTAMP NOT NULL,
	date_updated TIMESTAMP NOT NULL
);

CREATE UNIQUE INDEX filename_idx ON images (filename);
