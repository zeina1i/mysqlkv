package mysqlkv

import (
	"fmt"
	"testing"
	"time"
)

func TestKVService_PutThenGetThenDeleteThenGet(t *testing.T) {
	executor, err := NewSingleNodeExecutor(&MySQLConfig{
		Username: "mysqlkv_user",
		Password: "mysqlkv",
		Host:     "localhost",
		Port:     6035,
		DB:       "mysqlkv",
	})
	if err != nil {
		t.Error(err)
	}

	done := make(chan bool)
	gc := newGC(executor, 1*time.Minute, done)

	service := NewKVService(executor, gc)

	inputK := "hello"
	inputV := "bye"
	err = service.Put(inputK, inputV)
	if err != nil {
		t.Error(err)
	}

	outputV, err := service.Get(inputK)
	if err != nil {
		t.Error(err)
	}

	if inputV != outputV {
		t.Error(fmt.Sprintf("expected %s got %s", inputV, outputV))
	}

	err = service.Del(inputK)
	if err != nil {
		t.Error(err)
	}

	outputV, err = service.Get(inputK)
	if err.Error() != "sql: no rows in result set" {
		t.Error("expected notfound error but the record is available")
	}
}

func TestKVService_PutThenTTLThenGet(t *testing.T) {
	executor, err := NewSingleNodeExecutor(&MySQLConfig{
		Username: "mysqlkv_user",
		Password: "mysqlkv",
		Host:     "localhost",
		Port:     6035,
		DB:       "mysqlkv",
	})
	if err != nil {
		t.Error(err)
	}

	done := make(chan bool)
	gc := newGC(executor, 1*time.Minute, done)

	service := NewKVService(executor, gc)

	inputK := "hello"
	inputV := "bye"
	err = service.Put(inputK, inputV)
	if err != nil {
		t.Error(err)
	}

	err = service.TTL(inputK, time.Now().Unix())
	if err != nil {
		t.Error(err)
	}

	_, err = service.Get(inputK)
	if err != nil && err.Error() != "sql: no rows in result set" {
		t.Error("expected notfound error but the record is available")
	}
}
