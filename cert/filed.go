package cert

import (
	"os"
	"fmt"
	"crypto/tls"
)
//FiledCertificate contains all neccessary information to operate with.
//CertPath and KeyPath may be specified either by absolute or by relative path.
type FiledCertificate struct{
	CertPath string
	KeyPath string
}
//GetCertificate loads certificate and key directly from disk by a specified path.
//and then converts it to tls.Certificate array.
//Both files should be encrypted in PEM-block.
//If you would like to have built-in cedrtificate and key pair, you should use
//EmbeddedCertificate struct instead of this.
func (c *FiledCertificate) GetCertificate() ([]tls.Certificate, error) {
	if _, err := os.Stat(c.CertPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("Certificate file located at %v do not exist", c.CertPath)
	}
	if _, err := os.Stat(c.KeyPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("Key file located at %v do not exist", c.CertPath)
	}
	certFile, err := tls.LoadX509KeyPair(c.CertPath, c.KeyPath)
	if err != nil {
		return nil, err
	}
	return []tls.Certificate{certFile}, err
}