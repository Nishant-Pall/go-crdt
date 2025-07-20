package main

type LWWState[T any] struct {
	peer      string
	timestamp int32
	value     T
}

type LWWRegister[T any] struct {
	id    string
	state *LWWState[T]
}

func (reg *LWWRegister[T]) GetValue() T {
	return reg.state.value
}

func (reg *LWWRegister[T]) SetValue(value T) {

	reg.state = &LWWState[T]{peer: reg.state.peer, timestamp: reg.state.timestamp + 1, value: value}
}

func (reg *LWWRegister[T]) Merge(state *LWWState[T]) {

	remotePeer, remoteTimestamp := state.peer, state.timestamp
	localPeer, localTimestamp := reg.state.peer, reg.state.timestamp
	if localTimestamp > remoteTimestamp {
		return
	}
	if remoteTimestamp == localTimestamp && localPeer > remotePeer {
		return
	}

	reg.state = state
}
