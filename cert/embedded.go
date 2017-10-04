package cert

import (
	"encoding/hex"
	"crypto/aes"
	"crypto/tls"
)
//EmbeddedCertificate contains all neccessary information about Certificate. 
//In case if you would like to store your PEM-encoded certificate and key directly inside 
//application, you should use this type of Cerrtificate. 
type EmbeddedCertificate struct {
	CertString       string
	KeyString        string
	DecryptionSecret []byte
}
//GetCertificate decrypts embedded certificate string and key string to prepared tls.Certificate array.
//It is possible to use this type of struct without encryption. 
func (e *EmbeddedCertificate) GetCertificate() ([]tls.Certificate, error) {
	cert, key, err := e.decryptAll()
	if err != nil {
		return nil, err
	}
	pair, err := tls.X509KeyPair(cert, key)
	if err != nil {
		return nil, err
	}
	return []tls.Certificate{pair}, nil
}

func (e *EmbeddedCertificate) decryptAll() (cert, key []byte, err error) {
	cert, err = e.decryptCert()
	if err != nil {
		return nil, nil, err
	}
	key, err = e.decryptKey()
	if err != nil {
		return nil, nil, err
	}
	return cert, key, err
}

func (e *EmbeddedCertificate) decryptCert() ([]byte, error) {
	return decryptStringAes(e.CertString, e.DecryptionSecret)
}

func (e *EmbeddedCertificate) decryptKey() ([]byte, error) {
	return decryptStringAes(e.KeyString, e.DecryptionSecret)
}

func decryptStringAes(data string, secret []byte) ([]byte, error) {
	dataHex, err := hex.DecodeString(data)
	if err != nil {
		return nil, err
	}
	if len(secret) == 0 {
		return dataHex[:], nil
	}

	cipher, err := aes.NewCipher(secret)
	if err != nil {
		return nil, err
	}
	for i:=cipher.BlockSize(); i<len(dataHex); i+= cipher.BlockSize() {
		cipher.Decrypt(
			dataHex[i-cipher.BlockSize():i],
			dataHex[i-cipher.BlockSize():i])
	}
	return dataHex[:], nil
}