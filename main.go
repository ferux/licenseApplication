package main

//TODO: Add db drivers support.
import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/ferux/validationService/daemon"
	"github.com/ferux/validationService/db"
)

//Common info about application
var (
	Version = "v1.0.0"
	Author  = "Aleksandr Trushkin"
	Email   = "atrushkin@outlook.com"
)

const connString = "mongodb://%s:%s@%s/%s"

func main() {
	var dbUser, dbPassword, dbName, dbHost, dbCollection string
	var ver, auth, cont bool
	flag.StringVar(&dbUser, "user", "root", "Username for Database connection")
	flag.StringVar(&dbPassword, "pass", "root", "Password for Database")
	flag.StringVar(&dbName, "name", "database", "Name of database on the server")
	flag.StringVar(&dbCollection, "collection", "default", "Name of your collection")
	flag.StringVar(&dbHost, "host", "localhost", "Address of database")
	flag.BoolVar(&ver, "version", false, "Version of the application")
	flag.BoolVar(&auth, "author", false, "Author of the application")
	flag.BoolVar(&cont, "contact", false, "Author's contact info")
	flag.Parse()
	if ver {
		log.Println(Version)
	}
	if auth {
		log.Println(Author)
	}
	if cont {
		log.Println(Email)
	}
	cs := fmt.Sprintf(connString, dbUser, dbPassword, dbHost, dbName)
	log.Println(cs)
	dbConf := db.Config{
		Connection: cs,
		Database:   dbName,
		Collection: dbCollection,
	}
	go func() {
		_, err := daemon.Run(daemon.Config{DBConf: dbConf})
		if err != nil {
			log.Printf("Got an error while trying to run daemon. Reason:\n%v", err)
		}
	}()

	exitc := make(chan os.Signal)
	signal.Notify(exitc, syscall.SIGINT, syscall.SIGTERM)
	<-exitc
	log.Print("Got signal. Exiting!")
}

type estring string

func (s estring) ReplaceMe(old, new string) estring {
	return estring(strings.Replace(s.String(), old, new, -1))
}

func (s estring) String() string {
	return string(s)
}
