package page

import (
	"fmt"
	"net/http"
	"html/template"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"github.com/ezeev/saga/cloudsql"
	"os"
	"github.com/ezeev/saga/session"
)

func HandlePage(w http.ResponseWriter, r *http.Request,tmpl string, title string) {

	c := appengine.NewContext(r)

	globPattern := os.Getenv("TEMPLATE_GLOB_PATTERN")
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
}






