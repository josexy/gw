package handler

import "sync"

type MapBaseHandler[T any] struct {
	sync.RWMutex
	Cache map[string]T
	List  []T
}

func NewMapBaseHandler[T any]() MapBaseHandler[T] {
	return MapBaseHandler[T]{
		Cache: make(map[string]T),
	}
}
