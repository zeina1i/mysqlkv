package mysqlkv

import _ "github.com/go-sql-driver/mysql"

type Executor interface {
	initializeDB() error
	getKV(k string) (string, error)
	addOrUpdateKV(k string, v string) error
	updateKVExpiry(k string, expiry int64) error
	batchDeleteExpiredKVs(limit int) error
}
