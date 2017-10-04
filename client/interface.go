package client

import (
	"net/url"

	"github.com/ferux/dbgMe"
	"github.com/ferux/validationService/cert"
)

//IClient is an interface
type IClient interface {
	New(url.URL, cert.HandleCertificate, dbgMe.Debugger) *IClient
	RegisterMe() error
	Close() error
	RefreshClientInfo() error
}
