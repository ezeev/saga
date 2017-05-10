package page

import (
	"fmt"
	"net/http"
	"html/template"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"github.com/ezeev/saga/cloudsql"
	"github.com/ezeev/saga/session"
	"github.com/ezeev/saga/metrics"
	"github.com/ezeev/saga/config"
	"encoding/json"
)

func HandlePage(w http.ResponseWriter, r *http.Request, tmpl string, title string) {

	c := appengine.NewContext(r)
	conf, err := config.Config()
	if err != nil {
		log.Errorf(c,"There was an error loading conf", err)
	}

	globPattern := conf.TemplateGlobPattern
	if globPattern == "" {
		session.SetLastFailMsg(w,"TEMPLATE_GLOB_PATTERN env var not set!")
		log.Errorf(c,"TEMPLATE_GLOB_PATTERN env var not set!")
	}

	db, err := cloudsql.CloudSQLConnection()
	if err != nil {
		log.Errorf(c,"Error acquiring CloudSQL connection: %s", err.Error())
	}
	defer db.Close()
	pg, err := NewPage(w,r,db,"")
	if err != nil {

	}
	pg.Title = title
	t, err := template.New(tmpl).Funcs(PageFuncMap()).ParseGlob(globPattern)
	if err != nil {
		fmt.Println(err)
	}
	t.Execute(w,pg)
	registry := metrics.Registry()
	registry.IncRequests()

}

func HandleAccountApi(w http.ResponseWriter, r *http.Request) {


	//load the page
	db, err := cloudsql.CloudSQLConnection()
	if err != nil {
		panic(err)
	}
	jwtToken := r.Header.Get("X-Auth-Token")
	pg, err := NewPage(w,r,db,jwtToken)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(pg); err != nil {
		panic(err)
	}
}




