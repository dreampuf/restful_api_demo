package config

type RestfulAPIConfig struct {
	Host string
	Port uint

	DBHost, DBUser, DBPassword string
	DBPort uint
}
