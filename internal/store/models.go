package store

import (
	"time"
)

// Proxy ...
type Proxy struct {
	ID           int
	IPAddr       string
	Port         string
	HostURL      string
	TypeProxy    string
	Country      string
	PingCheck    *int
	DateCheck    *time.Time
	ItemDatetime time.Time
}
