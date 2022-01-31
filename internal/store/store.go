package store

// ProxyRepository ...
type ProxyRepository interface {
	Update(*Proxy) error
	Find(string) (*Proxy, error)
}
