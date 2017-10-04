package model

import (
	"time"
)
//License represents Model element for parsing.
type License struct {
	HostID      string    `json:"hostid,omitempty"`
	Status      int       `json:"status,omitempty"`
	Application string    `json:"application,omitempty"`
	Expiration  time.Time `json:"expiration,omitempty"`
}
//Constants for Status propose
const (
	StatusOk int = iota + 1
	StatusPending
	StatusBan
	StatusNew
)