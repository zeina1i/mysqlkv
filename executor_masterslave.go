package mysqlkv

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
)

type MasterSlaveExecutor struct {
	masterDB     *sqlx.DB
	slaveDBs     []*sqlx.DB
	masterConfig *MySQLConfig
	slavesConfig []*MySQLConfig
}

func NewMasterSlaveExecutor(masterConfig *MySQLConfig, slaveConfigs []*MySQLConfig) (*MasterSlaveExecutor, error) {
	dsn := masterConfig.Username + ":" + masterConfig.Password + "@" + "(" + masterConfig.Host + ":" + strconv.Itoa(masterConfig.Port) + ")/" + masterConfig.DB + "?parseTime=true"
	masterDB, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return nil, err
	}

	var slaveDBs []*sqlx.DB
	for _, slaveConfig := range slaveConfigs {
		dsn := slaveConfig.Username + ":" + slaveConfig.Password + "@" + "(" + slaveConfig.Host + ":" + strconv.Itoa(masterConfig.Port) + ")/" + slaveConfig.DB + "?parseTime=true"
		slaveDB, err := sqlx.Connect("mysql", dsn)
		if err != nil {
			return nil, err
		}

		slaveDBs = append(slaveDBs, slaveDB)
	}

	return &MasterSlaveExecutor{
		masterDB:     masterDB,
		slaveDBs:     slaveDBs,
		masterConfig: masterConfig,
		slavesConfig: slaveConfigs,
	}, nil
}

func (e *MasterSlaveExecutor) initializeDB() error {
	stmt := `
create table if not exists kvs
(
    k      varchar(128) not null,
    v      longtext     not null,
    expiry datetime     null,
    constraint kvs_k_uindex
        unique (k)
);

alter table kvs
    add primary key (k);
`
	_, err := e.masterDB.Exec(stmt)
	if err != nil {
		return err
	}

	return nil
}

func (e *MasterSlaveExecutor) getKV(k string) (string, error) {
	query := `
select k,v from kvs
where
      k = ?
      and (expiry > UNIX_TIMESTAMP() or expiry is null);
`
	var kv kv

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	slaveDBK := r1.Intn(len(e.slaveDBs) + 1)

	err := e.slaveDBs[slaveDBK].Get(&kv, query, k)
	if err != nil {
		return "", err
	}

	return kv.V, nil
}

func (e *MasterSlaveExecutor) addOrUpdateKV(k string, v string) error {
	query := `
replace into kvs
value (:k, :v, null)
`
	m := map[string]interface{}{
		"k": k,
		"v": v,
	}

	_, err := e.masterDB.NamedExec(query, m)

	return err
}

func (e *MasterSlaveExecutor) updateKVExpiry(k string, expiry int64) error {
	query := `
update kvs
set expiry = :expiry
where k = :k
`
	m := map[string]interface{}{
		"expiry": expiry,
		"k":      k,
	}

	_, err := e.masterDB.NamedExec(query, m)

	return err
}

func (e *MasterSlaveExecutor) batchDeleteExpiredKVs(limit int) error {
	query := `
delete from kvs
where expiry < UNIX_TIMESTAMP()
limit :limit
`
	m := map[string]interface{}{
		"limit": limit,
	}

	_, err := e.masterDB.NamedExec(query, m)

	return err
}
