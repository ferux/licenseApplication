package db

import (
	"fmt"
)

//Config contains connection info for db.
type Config struct {
	Connection string
	User string
	Password string
	Server string
	Port string
	Database string
	Collection string
}
//PrepareConfig uses fields like User, Password, Server, Port and composites it to Connection string.
func (c *Config) PrepareConfig() {
	c.Connection = fmt.Sprintf("mongo://%s:%s@%s:%s/%s", c.User, c.Password, c.Server, c.Port, c.Database)
}