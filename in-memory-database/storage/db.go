package storage

import "errors"

var (
	ErrorNotFound = errors.New("key not found")
)

type Db interface {
	Get(key string) (string, error)
	Set(key string, val string) (string, error)
	Del(key string) (string, error)
}
