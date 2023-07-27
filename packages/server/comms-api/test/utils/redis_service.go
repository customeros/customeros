package utils

import "context"

type DatabaseValues struct {
	Tenant string
	Active bool
}

type TestRedisService struct {
	KeyMap map[string]DatabaseValues
}

func (s *TestRedisService) GetKeyInfo(ctx context.Context, tag, key string) (bool, *string) {
	if val, ok := s.KeyMap[key]; ok {
		return val.Active, &val.Tenant
	} else {
		return false, nil
	}
}

func NewTestRedisService() *TestRedisService {
	return &TestRedisService{
		KeyMap: make(map[string]DatabaseValues),
	}
}
