package mysqlkv

type KVService struct {
	store Store
	gc    *GC
}

func NewKVService(store Store, gc *GC) *KVService {
	gc.collect()

	return &KVService{store: store}
}

func (s *KVService) Put(k string, v string) error {
	err := s.store.addOrUpdateKV(k, v)

	return err
}

func (s *KVService) Get(k string) (string, error) {
	k, err := s.store.getKV(k)
	if err != nil {
		return "", err
	}

	return k, err
}

func (s *KVService) Del(k string) error {
	err := s.store.updateKVExpiry(k, -1)

	return err
}

func (s *KVService) TTL(k string, expiry int64) error {
	err := s.store.updateKVExpiry(k, expiry)

	return err
}
