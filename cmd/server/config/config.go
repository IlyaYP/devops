package config

// move defaults to flag Parse
//type Config struct {
//	Address       string `env:"ADDRESS" envDefault:"localhost:8080"`
//	StoreInterval int    `env:"STORE_INTERVAL" envDefault:"300"`
//	StoreFile     string `env:"STORE_FILE" envDefault:"/tmp/devops-metrics-db.json"`
//	Restore       bool   `env:"RESTORE" envDefault:"true"`
//}

type Config struct {
	Address       string `env:"ADDRESS"`
	StoreInterval int    `env:"STORE_INTERVAL"`
	StoreFile     string `env:"STORE_FILE"`
	Restore       bool   `env:"RESTORE"`
}
