package main

import (
	"config"
	"flag"
	"db"
	"fmt"
	"os"
)

func main() {
	cfg := &config.RestfulAPIConfig{}

	flag.StringVar(&cfg.DBHost, "dbhost", "127.0.0.1", "Connecting DB Host")
	flag.UintVar(&cfg.DBPort, "dbport", 5432, "Connecting DB Port")
	flag.StringVar(&cfg.DBUser, "dbuser", "api", "Connecting DB User")
	// CRITICAL: Need more security here
	flag.StringVar(&cfg.DBPassword, "dbpassword", "api", "Connecting DB Password")
	flag.Parse()

	apiDataSource := db.NewAPIDataSource(cfg)

	err := db.MigrateRun(apiDataSource)
	if err != nil {
		fmt.Fprintf(os.Stderr, fmt.Sprintf("Error: %s\n", err.Error()))
		os.Exit(1)
	}
}
