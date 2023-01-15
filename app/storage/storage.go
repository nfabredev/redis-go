package storage

import "errors"

type Storage map[string]string

var storage Storage

func init() {
	storage = make(Storage)
}

func Set(key string, value string) Storage {
	storage[key] = value
	return storage
}

func Get(key string) (string, error) {
	value, ok := storage[key]
	if !ok {
		return "", errors.New("-ERR get called on unset key: " + key)
	}
	return value, nil
}
