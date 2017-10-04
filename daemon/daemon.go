//Package daemon serves api for webservice.
//Usage: the api is following:
// To show all available licenses just send an empty get request to "/"
// To show specified license visit path /read/{:id}
// To create | update (maybe i will split those two functions later) license use
// post method with the following parameters:
//		hostid -- license (required)
//		status -- license status (optional)
//		application -- app name (optional)
//		expires -- date of expire (optional)
// and send it to "/update"
//TODO:
// Authentication check
// Delete method
// Split Create and Update
//
package daemon

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ferux/validationService/db"
	"github.com/ferux/validationService/model"
)

//Config represents a config for daemons
type Config struct {
	DBConf db.Config
}

//DB represents M in MVC-framework
type DB struct {
	model *model.Model
}

//Run a server.
func Run(conf Config) (*model.Model, error) {
	session, err := db.New(conf.DBConf)
	if err != nil {
		return nil, err
	}
	m := model.New(session)
	http.Handle("/", handleLicenses(m))
	http.Handle("/read/", handleGetLicense(m))
	http.Handle("/update/", handleInsertLicense(m))
	go func() {
		log.Print(http.ListenAndServe("127.0.0.1:6666", nil))
	}()
	return m, nil

}

func handleLicenses(m *model.Model) http.Handler {
	listAllFunc := func(w http.ResponseWriter, r *http.Request) {
		licenses, err := m.SelectLicenses()
		if err != nil {
			http.Error(w, "Got an error\n"+err.Error(), http.StatusBadRequest)
			return
		}
		js, err := json.Marshal(&licenses)
		if err != nil {
			http.Error(w, "Got an error\n"+err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Fprint(w, string(js))
	}

	return http.HandlerFunc(listAllFunc)
}

func handleGetLicense(m *model.Model) http.Handler {
	getFunc := func(w http.ResponseWriter, r *http.Request) {
		getID := strings.Split(r.URL.Path, "/")[2]
		log.Println("Finding id:", getID)
		license, err := m.SelectLicense(getID)
		if err != nil {
			http.Error(w, "Got an error\n"+err.Error(), http.StatusInternalServerError)
			return
		}
		output, err := json.Marshal(&license)
		if err != nil {
			http.Error(w, "Got an error\n"+err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprint(w, string(output))
	}
	return http.HandlerFunc(getFunc)
}

func handleInsertLicense(m *model.Model) http.Handler {
	insertFunc := func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		if r.Method != "POST" {
			http.Error(w, "{\"status\":\"Method should be post\"}", http.StatusInternalServerError)
			return
		}
		st, err := strconv.Atoi(r.FormValue("status"))
		if err != nil {
			st = model.StatusBan
		}
		dt, err := time.Parse("2006-01-02", r.FormValue("expiration"))
		if err != nil {
			dt = time.Now()
		}
		app := r.FormValue("application")
		license := &model.License{
			HostID:      r.FormValue("hostid"),
			Status:      st,
			Expiration:  dt,
			Application: app,
		}
		err = m.UpdateLicense(license)
		if err != nil {
			http.Error(w, "{\"status\":\""+err.Error()+"\"}", http.StatusInternalServerError)
		}
		w.Header().Add("Status Code", "200")
		fmt.Fprintf(w, `{"status":"ok"}`)
	}
	return http.HandlerFunc(insertFunc)
}
