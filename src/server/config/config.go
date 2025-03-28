package config

type Config struct {
	HTTP      `json:"http"`
	Databases struct {
		Shop     Database `json:"shop"`
		Images   Database `json:"images"`
		Printful Database `json:"printful"`
	} `json:"databases"`
	Printful `json:"printful"`
	Images   `json:"images"`
	Sessions `json:"sessions"`
	Paypal   `json:"paypal"`
}

type HTTP struct {
	Port          int    `json:"port"`
	HttpsKeyFile  string `json:"https_key_file"`
	HttpsCertFile string `json:"https_cert_file"`
}

type Database struct {
	ConnectURI string   `json:"connect_uri"`
	DBName     string   `json:"db_name"`
	BucketName string   `json:"bucket_name"`
	KeyVault   KeyVault `json:"key_vault"`
}

type Printful struct {
	Endpoint        string  `json:"endpoint"`
	AccessToken     string  `json:"access_token"`
	SimulateMockup  bool    `json:"simulate_mockup"`
	SimulateTaskKey string  `json:"simulate_task_key"`
	TaskInterval    int     `json:"task_interval"`
	MockupDirectory string  `json:"mockup_directory"`
	ImagesURL       string  `json:"images_url"`
	Markup          float64 `json:"markup"`
}

type Images struct {
	BaseURL string `json:"base_url"`
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

type KeyVault struct {
	KMS        KMS    `json:"kms"`
	DBName     string `json:"db_name"`
	Collection string `json:"collection"`
	DEKID      string `json:"dek_id"`
}

type KMS struct {
	Endpoint        string `json:"endpoint"`
	CertificatePath string `json:"certificate_path"`
}
