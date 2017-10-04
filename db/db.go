package db

import (
	"fmt"
	"log"

	"github.com/ferux/validationService/model"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//Client is a struct that contains all main information about connection to db
type Client struct {
	Conn       *mgo.Session
	Collection *mgo.Collection
	Config     Config
	Close      func()
}

//New opens connection to db and creates a struct with information
func New(c Config) (*Client, error) {
	if len(c.Collection) == 0 {
		c.PrepareConfig()
	}
	if len(c.Connection) < len("mongo://a:a@a:1/") {
		return nil, fmt.Errorf("Your connectionstring is too small: %s", c.Connection)
	}
	conn, err := mgo.Dial(c.Connection)
	if err != nil {
		return nil, err
	}
	if err := conn.Ping(); err != nil {
		conn.Clone()
		return nil, err
	}
	log.Print("Ping to database was successful!")
	conn.SetMode(mgo.Monotonic, true)
	collection := conn.DB(c.Database).C(c.Collection)
	if collection == nil {
		log.Print("Can't find collection.")
	}
	closeFunc := func() { conn.Close() }
	return &Client{Conn: conn, Config: c, Close: closeFunc, Collection: collection}, nil
}

//SelectLicenses returns an array of all Licenses.
func (c *Client) SelectLicenses() ([]*model.License, error) {
	result := make([]*model.License, 1)
	err := c.Collection.Find(bson.M{}).All(&result)
	return result, err
}

//SelectLicense returns a specified license.
func (c *Client) SelectLicense(license string) (*model.License, error) {
	var result *model.License
	err := c.Collection.Find(bson.M{"hostid": license}).One(&result)
	// if err != nil {log.Printf("Got error while finding license: %v", err)}
	return result, err
}

//UpdateLicense updates specified license.
func (c *Client) UpdateLicense(l *model.License) error {
	if _, err := c.SelectLicense(l.HostID); err != nil && err.Error() == "not found" {
		return c.InsertLicense(l)
	} else if err != nil {
		return err
	}
	return c.Collection.Update(bson.M{"hostid": l.HostID}, l)
}

//InsertLicense inserts a row to db
func (c *Client) InsertLicense(l *model.License) error {
	return c.Collection.Insert(l)
}
