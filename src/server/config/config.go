package config

type Config struct {
	HTTP     HTTP     `json:"http"`
	Database Database `json:"database"`
	Printful Printful `json:"printful"`
}

type HTTP struct {
	Port          int    `json:"port"`
	HttpsKeyFile  string `json:"https_key_file"`
	HttpsCertFile string `json:"https_cert_file"`
}

type Database struct {
	ConnectURI string `json:"connect_uri"`
	DBName     string `json:"db_name"`
}

type Printful struct {
	Endpoint string `json:"endpoint"`
}
