// Package mail provides an HTTP endpoint for sending mail. This package is currently under development.
package mail


import (
	"fmt"
	"net/http"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/mail"
	"os"
	"github.com/ezeev/saga/session"
)


//NOT IN ACTIVE USE

func RegisterHandlers() {

	path := "/contact"
	envpath := os.Getenv("MAIL_CONTACTFORM_HANDLER")
	if envpath != "" {
		path = envpath
	}

	http.HandleFunc(path, HandleContactForm)
}

func HandleContactForm(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	addr := r.FormValue("email")
	name := r.FormValue("name")
	content := r.FormValue("message")

	sender := os.Getenv("MAIL_SENDER_ADDRESS")
	receiver := os.Getenv("MAIL_RECEIVER_ADDRESS")
	subject := os.Getenv("MAIL_CONTACTFORM_SUBJECT")
	redirect := os.Getenv("MAIL_CONTACTFORM_REDIRECT")

	msg := &mail.Message{
		Sender:  sender,
		To:      []string{receiver},
		Subject: subject,
		Body:    fmt.Sprintf(confirmMessage,name,addr,content),
	}
	if err := mail.Send(ctx, msg); err != nil {
		log.Errorf(ctx, "Couldn't send email: %v", err)
	}
	session.SetLastSuccessMsg(w, "Your message has been sent.")
	http.Redirect(w, r, redirect, http.StatusTemporaryRedirect)
}

const confirmMessage = `
The messasge below is from the contact form:

Name: %s

Email: %s

Message (below this line)
 ------------------------------------------
 %s
`