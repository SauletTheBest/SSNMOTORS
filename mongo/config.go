package mongo

import "fmt"

type Config struct {
	Database     string `env:"MONGO_DB"`
	URI          string `env:"MONGO_DB_URI"`
	Username     string `env:"MONGO_USERNAME"`
	Password     string `env:"MONGO_PWD"`
	ReplicaSet   string `env:"MONGO_DB_REPLICA_SET"`
	WriteConcern string `env:"MONGO_WRITE_CONCERN"`
	TLSFilePath  string `env:"MONGO_TLS_FILE_PATH"`
	TLSEnable    bool   `env:"MONGO_TLS_ENABLE" envDefault:"false"`
}

func (m Config) genConnectURL() string {
	var url string
	if m.Username == "" || m.Password == "" {
		url = fmt.Sprintf("mongodb://%s/?tls=%t&retryWrites=false", m.URI, m.TLSEnable)
	} else {
		url = fmt.Sprintf("mongodb://%s:%s@%s/?tls=%t&retryWrites=false", m.Username, m.Password, m.URI, m.TLSEnable)
	}

	return url
}
