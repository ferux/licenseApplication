package cert

import (
	"crypto/tls"
)
//HandleCertificate ensures type can return tls.Certificate array.
type HandleCertificate interface {
	GetCertificate() ([]tls.Certificate, error)
}