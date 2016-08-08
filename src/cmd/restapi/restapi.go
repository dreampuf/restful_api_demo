package main

import (
	"config"
	"flag"
	"service"
	"db"
)

func main() {
	cfg := &config.RestfulAPIConfig{}

	flag.StringVar(&cfg.Host, "host", "127.0.0.1", "WebService Host")
	flag.UintVar(&cfg.Port, "port", 8080, "WebService Port")
	flag.StringVar(&cfg.DBHost, "dbhost", "127.0.0.1", "Connecting DB Host")
	flag.UintVar(&cfg.DBPort, "dbport", 5432, "Connecting DB Port")
	flag.StringVar(&cfg.DBName, "dbname", "api", "Connecting DB Name")
	flag.StringVar(&cfg.DBUser, "dbuser", "api", "Connecting DB User")
	// CRITICAL: Need more security here
	flag.StringVar(&cfg.DBPassword, "dbpassword", "api", "Connecting DB Password")
	flag.Parse()

	apiDataSource := db.NewAPIDataSource(cfg)
	defer apiDataSource.Close()
	web := service.NewWebService(cfg.Host, cfg.Port, apiDataSource)
	web.Serve()
}
