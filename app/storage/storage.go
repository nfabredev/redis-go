package storage

import "errors"

type Value struct {
	Value string
	Exp   int64
}

type Storage map[string]Value

var storage Storage

func init() {
	storage = make(Storage)
}

func Set(key string, value Value) Storage {
	storage[key] = value
	return storage
}

func Get(key string) (Value, error) {
	value, ok := storage[key]
	if !ok {
		return Value{}, errors.New("-ERR get called on unset key: " + key)
	}
	return value, nil
}
