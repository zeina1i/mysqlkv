package mysqlkv

type kv struct {
	K      string `db:"k"`
	V      string `db:"v"`
	Expiry int64  `db:"expiry"`
}
