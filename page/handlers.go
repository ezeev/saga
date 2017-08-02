package page

import (
	"net/http"
	"html/template"
	"log"
	"github.com/ezeev/saga/cloudsql"
	"github.com/ezeev/saga/session"
	"github.com/ezeev/saga/metrics"
	"github.com/ezeev/saga/config"
	"encoding/json"
)

func HandlePage(w http.ResponseWriter, r *http.Request, tmpl string, title string) {

	conf, err := config.Config()
	if err != nil {
		log.Printf("There was an error loading conf, %s", err)
	}

	globPattern := conf.TemplateGlobPattern
	if globPattern == "" {
		session.SetLastFailMsg(w,"TEMPLATE_GLOB_PATTERN env var not set!")
		log.Print("TEMPLATE_GLOB_PATTERN env var not set!")
	}

	db, err := cloudsql.CloudSQLConnection()
	if err != nil {
		log.Printf("Error acquiring CloudSQL connection: %s", err.Error())
	}
	defer db.Close()
	pg, err := NewPage(w,r,db,"")
	if err != nil {
		log.Printf("Error creating Page instance: %s", err)
	}
	pg.Title = title

	t, err := template.New(tmpl).Funcs(PageFuncMap()).ParseGlob(globPattern)
	if err != nil {
		log.Print(err)
	}
	w.WriteHeader(http.StatusOK)
	t.Execute(w,pg)
	registry := metrics.Registry()
	registry.IncRequests()

}

func HandlePageApi(w http.ResponseWriter, r *http.Request) {


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




