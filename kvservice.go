package mysqlkv

import (
	"log"
)

type KVService struct {
	store Store
}

func NewKVService(store Store) *KVService {
	return &KVService{store: store}
}

func (s *KVService) Put(k string, v string) error {
	err := s.store.addOrUpdateKV(k, v)

	return err
}

func (s *KVService) Get(k string) (string, error) {
	k, err := s.store.getKV(k)
	if err != nil {
		//log.Fatal(err)
		return "", err
	}

	return k, err
}

func (s *KVService) Del(k string) error {
	err := s.store.updateKVExpiry(k, -1)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func (s *KVService) TTL(k string, expiry int64) error {
	err := s.store.updateKVExpiry(k, expiry)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return err
}
