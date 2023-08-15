package mysqlkv

type kv struct {
	K      string `db:"k"`
	V      string `db:"v"`
	Expiry int64  `db:"expiry"`
}

type MySQLConfig struct {
	Username string
	Password string
	Host     string
	Port     int
	DB       string
}
