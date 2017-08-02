package mailgun

import (
	"log"
	"net/http"
	"github.com/ezeev/saga/session"
	"github.com/ezeev/saga/config"
	"gopkg.in/mailgun/mailgun-go.v1"
)

var conf *config.SagaConfig

func init() {

	log.Print("Regisering Mailgun Handlers:")
	http.HandleFunc("/sendmail", HandleSendMail)
	log.Print("/sendmail")

	conf, _ = config.Config()
}

/*
curl https://viroonga.com/sendmail \
	-F from='Website Inbound Mail <donotreply@viroonga.com>' \
	-F subject='Hello Test' \
	-F body='Testing some Mailgun awesomness!'
 */


func HandleSendMail(w http.ResponseWriter, r *http.Request) {

	log.Print("Received mail request")

	referrer := r.Referer()

	recip := conf.MailDefaultRecipientEmail
	from := r.PostFormValue("from")
	name := r.PostFormValue("name")
	//subject := r.PostFormValue("subject")
	subject := "Website mail"
	var body string
	if name != "" {
		body = "Message from " + name + ":\n\n"
	}
	body = body + r.PostFormValue("body")

	mg := mailgun.NewMailgun(conf.MailGunDomain, conf.MailGunApiKey, conf.MailGunPubKey)
	message := mg.NewMessage(from,subject,body,recip)
	_, id, err := mg.Send(message)
	if err != nil {
		log.Print("Error sending mailgun message: %s", err)
	} else {
		log.Print("Sent message id: %s", id)
	}

	session.SetLastSuccessMsg(w,"Your message has been sent!")
	http.Redirect(w, r, referrer, http.StatusTemporaryRedirect)
}

