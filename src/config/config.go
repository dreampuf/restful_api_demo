package config

type RestfulAPIConfig struct {
	Host string
	Port uint

	DBHost, DBName, DBUser, DBPassword string
	DBPort uint
}
