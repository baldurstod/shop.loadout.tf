package config

type Config struct {
	HTTP      `json:"http"`
	Databases struct {
		Shop   Database `json:"shop"`
		Images Database `json:"images"`
	} `json:"databases"`
	Printful `json:"printful"`
	Sessions `json:"sessions"`
	Paypal   `json:"paypal"`
}

type HTTP struct {
	Port          int    `json:"port"`
	HttpsKeyFile  string `json:"https_key_file"`
	HttpsCertFile string `json:"https_cert_file"`
}

type Database struct {
	ConnectURI string `json:"connect_uri"`
	DBName     string `json:"db_name"`
	BucketName string `json:"bucket_name"`
}

type Printful struct {
	Endpoint string `json:"endpoint"`
}

type Sessions struct {
	ConnectURI  string `json:"connect_uri"`
	DBName      string `json:"db_name"`
	Collection  string `json:"collection"`
	Secret      string `json:"secret"`
	SessionName string `json:"session_name"`
}

type Paypal struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}
