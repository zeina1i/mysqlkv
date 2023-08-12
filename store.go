package mysqlkv

type Store interface {
	getKV(k string) (string, error)
	addOrUpdateKV(k string, v string) error
	updateKVExpiry(k string, expiry int64) error
	//batchDeleteExpiredKVs(num int)
}
