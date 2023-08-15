package mysqlkv

type KVService struct {
	executor Executor
	gc       *GC
}

func NewKVService(executor Executor, gc *GC) *KVService {
	gc.collect()

	return &KVService{executor: executor}
}

func (s *KVService) Put(k string, v string) error {
	err := s.executor.addOrUpdateKV(k, v)

	return err
}

func (s *KVService) Get(k string) (string, error) {
	k, err := s.executor.getKV(k)
	if err != nil {
		return "", err
	}

	return k, err
}

func (s *KVService) Del(k string) error {
	err := s.executor.updateKVExpiry(k, -1)

	return err
}

func (s *KVService) TTL(k string, expiry int64) error {
	err := s.executor.updateKVExpiry(k, expiry)

	return err
}
