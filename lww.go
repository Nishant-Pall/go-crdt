package main

import "maps"

type LWWMapState[T any] map[string]LWWRegister[T]

type LWWMap[T any] struct {
	id   string
	data LWWMapState[T]
}

func (lwwMap *LWWMap[T]) DefineLWWMap(id string, state LWWMapState[T]) {
	lwwMap.id = id
	maps.Copy(state, lwwMap.data)
}

func (lwwMap *LWWMap[T]) GetValue() LWWMapState[T] {
	return lwwMap.data
}

func (lwwMap *LWWMap[T]) Merge(state LWWMapState[T]) {
	for key, remote := range state {

		local, ok := lwwMap.data[key]

		if !ok {
			local.Merge(remote.state)
		} else {
			lwwMap.data[key] = LWWRegister[T]{id: lwwMap.id, state: remote.state}
		}
	}
}

func (lwwMap *LWWMap[T]) SetValue(key string, value T) {
	register, ok := lwwMap.data[key]
	if !ok {
		register.SetValue(value)
	} else {
		state := &LWWState[T]{peer: lwwMap.id, timestamp: 1, value: value}
		lwwMap.data[key] = LWWRegister[T]{id: lwwMap.id, state: state}
	}

}

func (lwwMap *LWWMap[T]) Get(key string) T {
	return lwwMap.data[key].state.value
}

func (lwwMap *LWWMap[T]) Delete(key string) {
	lwwMap.data[key] = LWWRegister[T]{}
}

func (lwwMap *LWWMap[T]) Has(key string) bool {
	_, ok := lwwMap.data[key]
	return !ok
}
