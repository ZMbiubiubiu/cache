package cache

type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}

type PeerGetter interface {
	PeerGet(group string, key string) ([]byte, error)
}
