package redis

import (
	"time"

	r "github.com/go-redis/redis"
)

type Storage struct {
	client *r.Client
}

func NewStorage(url string) (*Storage, error) {
	var opts *r.Options
	var err error

	if opts, err = r.ParseURL(url); err != nil {
		return nil, err
	}

	return &Storage{
		client: r.NewClient(opts),
	}, nil
}

func (s *Storage) Get(key string) []byte {
	val, _ := s.client.Get(key).Bytes()
	return val
}

func (s *Storage) Set(key string, content []byte, duration time.Duration) {
	s.client.Set(key, content, duration)
}
