package config

type Config struct {
	HTTP     HTTP     `json:"http"`
	Database Database `json:"database"`
	Printful Printful `json:"printful"`
	Sessions Sessions `json:"sessions"`
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
}

type Printful struct {
	Endpoint string `json:"endpoint"`
}

type Sessions struct {
	Path       string `json:"path"`
	AuthKey    string `json:"auth_key"`
	EncryptKey string `json:"encrypt_key"`
}

type Paypal struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}
