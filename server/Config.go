package server

import (
	"crypto/tls"

	"github.com/ferux/validationService/model"
)

//Config represents the neccessary parameters for running the server.
type Config struct {
	IP              string
	Port            string
	ServerCertPath  string
	ClientCertsPath string
	ServerConfig    *tls.Config
	Model           *model.Model
}
