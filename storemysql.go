package mysqlkv

import (
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type MySQLConfig struct {
	Username string
	Password string
	Host     string
	Port     int
	DB       string
}

type StoreMySQL struct {
	db     *sqlx.DB
	config *MySQLConfig
}

func NewStoreMysql(config *MySQLConfig) (*StoreMySQL, error) {
	dsn := config.Username + ":" + config.Password + "@" + "(" + config.Host + ":" + strconv.Itoa(config.Port) + ")/" + config.DB + "?parseTime=true"
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {

		return nil, err
	}

	return &StoreMySQL{
		db:     db,
		config: config,
	}, nil
}

func (s *StoreMySQL) InitializeDB() error {
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
	_, err := s.db.Exec(stmt)
	if err != nil {
		return err
	}

	return nil
}

func (s *StoreMySQL) getKV(k string) (string, error) {
	query := `
select k,v from kvs
where
      k = ?
      and (expiry > UNIX_TIMESTAMP() or expiry is null);
`
	var kv kv

	err := s.db.Get(&kv, query, k)
	if err != nil {
		return "", err
	}

	return kv.V, nil
}

func (s *StoreMySQL) addOrUpdateKV(k string, v string) error {
	query := `
replace into kvs
value (:k, :v, null)
`
	m := map[string]interface{}{
		"k": k,
		"v": v,
	}

	_, err := s.db.NamedExec(query, m)

	return err
}

func (s *StoreMySQL) updateKVExpiry(k string, expiry int64) error {
	query := `
update kvs
set expiry = :expiry
where k = :k
`
	m := map[string]interface{}{
		"expiry": expiry,
		"k":      k,
	}

	_, err := s.db.NamedExec(query, m)

	return err
}

func (s *StoreMySQL) batchDeleteExpiredKVs(limit int) error {
	query := `
delete from kvs
where expiry < UNIX_TIMESTAMP()
limit :limit
`
	m := map[string]interface{}{
		"limit": limit,
	}

	_, err := s.db.NamedExec(query, m)

	return err
}
