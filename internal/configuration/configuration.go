package configuration

type StConfig struct {
	Host     string
	Port     string
	Database string
	Username string
	Password string
}

var StorageConfig = StConfig{
	"localhost",
	"5432",
	"wblab0",
	"postgres",
	"password"}
