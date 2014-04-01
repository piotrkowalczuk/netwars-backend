package database


type RedisConfig struct {
	Address string `xml:"address"`
}

type PostgreConfig struct {
	ConnectionString string `xml:"connectionString"`
}
