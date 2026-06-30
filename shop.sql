--enable encryption
CREATE EXTENSION pgcrypto ;

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
	id BIGSERIAL PRIMARY KEY,
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

CREATE TABLE contacts (
	id BIGSERIAL PRIMARY KEY,
	subject TEXT NOT NULL,
	email TEXT NOT NULL,
	content TEXT NOT NULL,
	status TEXT NOT NULL,
	date_created TIMESTAMP NOT NULL
);

--CREATE TYPE address AS (
--	first_name TEXT,
--	last_name TEXT,
--	organization TEXT,
--	address1 TEXT,
--	address2 TEXT,
--	city TEXT,
--	state_code TEXT,
--	state_name TEXT,
--	country_code TEXT,
--	country_name TEXT,
--	postal_code TEXT,
--	phone TEXT,
--	email TEXT,
--	tax_number TEXT
--);

--CREATE TYPE order_item AS (
--	product_id TEXT,
--	name TEXT,
--	product JSONB,
--	quantity INTEGER,
--	retail_price DECIMAL,
--	thumbnail_url TEXT
--);

CREATE TYPE cart AS (
	currency TEXT,
	items TEXT[]
);

CREATE TABLE orders (
	id TEXT PRIMARY KEY,
	currency TEXT NOT NULL,
	shipping_address BYTEA NOT NULL,
	billing_address BYTEA,
	same_billing_address BOOLEAN NOT NULL,
	items JSONB NOT NULL,
	shipping_infos JSONB NOT NULL,
	tax_info JSONB NOT NULL,
	shipping_method TEXT NOT NULL,
	printful_order_id TEXT NOT NULL,
	paypal_order_id TEXT NOT NULL,
	dek BYTEA NOT NULL,
	status TEXT NOT NULL,
	date_created TIMESTAMP NOT NULL,
	date_updated TIMESTAMP NOT NULL
);

CREATE TABLE users (
	id TEXT PRIMARY KEY,
	username TEXT NOT NULL,
	password TEXT NOT NULL,
	display_name TEXT NOT NULL,
	email_verified BOOLEAN NOT NULL,
	address BYTEA NOT NULL,
	currency TEXT NOT NULL,
	orders TEXT[] NOT NULL,
	favorites TEXT[] NOT NULL,
	cart JSONB NOT NULL,
	dek BYTEA NOT NULL,
	date_created TIMESTAMP NOT NULL,
	date_updated TIMESTAMP NOT NULL
);
