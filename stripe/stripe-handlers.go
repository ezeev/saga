package stripe

import (

	"net/http"
	"github.com/ezeev/saga/cloudsql"
	"github.com/ezeev/saga/session"
	"github.com/ezeev/saga/util"
	"net/url"
	"log"
	"github.com/ezeev/saga/config"
)

var reqVars = []string{"STRIPE_TEST_PK", "STRIPE_TEST_SK", "STRIPE_LIVE_PK", "STRIPE_LIVE_SK"}


// RegisterHandlers registers the /addcard, /delcard, /subscribe, and /unsubscribe API endpoints.
func init() {

	err := util.CheckVars(reqVars)
	if err != nil {
		panic(err)
	}
	log.Print("Regisering Stripe Handlers:")
	http.HandleFunc("/addcard", HandleAddCard)
	http.HandleFunc("/delcard", HandleRemoveCard)
	http.HandleFunc("/subscribe", HandleSubscribe)
	http.HandleFunc("/unsubscribe", HandleUnsubscripe)
	log.Print("/addcard, /delcard, /subscribe, /unsubscribe")
}

func HandleAddCard(w http.ResponseWriter, r *http.Request) {

	conf, err := config.Config()
	if err != nil {
		log.Printf("Error loading config %s", err)
	}

	db, _ := cloudsql.CloudSQLConnection()
	defer db.Close()
	stripeMgr, err := NewStripeMgr(db)
	if err != nil {
		log.Printf("Error creating stipe mgr %s", err)
	}

	prof := session.Profile(r)
	email := prof.Email
	cust, err := stripeMgr.Customer(email)
	if err != nil {
		errS := "Error retreiving customer: %s"
		log.Printf(errS, err)
	}

	token := r.FormValue("token_id")
	log.Printf("All Form: %s", r.Form)
	log.Printf("Token: %s", token)
	log.Printf("Stripe Id: %s\n Email: %s\nToken ID:%s", cust.ID, email, token)

	//add the CC to the customer's acct
	card, err := stripeMgr.AddCard(cust.ID,token)
	//redir := os.Getenv("STRIPE_CARD_REDIRECT")
	redir := conf.StripeCardRedirect
	if session.LastReferrerUrl(r) != "" {
		redir = session.LastReferrerUrl(r)
		redir, _ = url.QueryUnescape(redir)
	}

	if err != nil {
		session.SetLastFailMsg(w,"Unable to add new card: " + err.Error())
	} else {
		session.SetLastSuccessMsg(w,"Successfully added card ending in " + card.LastFour + "!")
	}
	http.Redirect(w, r, redir, http.StatusTemporaryRedirect)
}

func HandleRemoveCard(w http.ResponseWriter, r *http.Request) {

	conf, err := config.Config()
	if err != nil {
		log.Printf("Error loading config %s", err)
	}

	db, _ := cloudsql.CloudSQLConnection()
	defer db.Close()
	stripeMgr, err := NewStripeMgr(db)
	if err != nil {
		errS := "Error creating stripemgr: %s"
		log.Printf(errS, err)

	}

	prof := session.Profile(r)
	email := prof.Email
	cust, err := stripeMgr.Customer(email)
	if err != nil {
		errS := "Error retreiving customer: %s"
		log.Printf(errS, err)
	}
	custId := cust.ID
	cardId := r.FormValue("cardId")
	log.Printf("Deleteing card for customer %s, %s", email, custId)

	_, err = stripeMgr.DelCard(custId,cardId)
	if err != nil {
		session.SetLastFailMsg(w, "Unable to delete card. " + err.Error())
	} else {
		session.SetLastSuccessMsg(w, "Successfully removed card!")
	}

	//redirect := os.Getenv("STRIPE_CARD_REDIRECT")
	redirect := conf.StripeCardRedirect
	if session.LastReferrerUrl(r) != "" {
		redirect = session.LastReferrerUrl(r)
		redirect, _ = url.QueryUnescape(redirect)
	}
	http.Redirect(w, r, redirect, http.StatusTemporaryRedirect)
}


func HandleSubscribe(w http.ResponseWriter, r *http.Request) {


	conf, err := config.Config()
	if err != nil {
		log.Printf("Error loading config %s", err)
	}

	db, _ := cloudsql.CloudSQLConnection()
	defer db.Close()
	stripeMgr, err := NewStripeMgr(db)

	redirect := conf.StripeCardRedirect //os.Getenv("STRIPE_CARD_REDIRECT")
	if session.LastReferrerUrl(r) != "" {
		redirect = session.LastReferrerUrl(r)
		redirect, _ = url.QueryUnescape(redirect)
	}

	if err != nil {
		errS := "Error creating stripemgr: %s"
		log.Printf(errS, err)
		session.SetLastFailMsg(w,"Unable to subscripe: " + err.Error())
		http.Redirect(w, r, redirect, http.StatusTemporaryRedirect)
	}

	prof := session.Profile(r)
	email := prof.Email
	planId := r.URL.Query().Get("planId")
	if planId == "" {
		panic("planId is required")
	}
	_, err = stripeMgr.SubscribeCustomer(email,planId)
	if err != nil {
		session.SetLastFailMsg(w,"Unable to subscripe: " + err.Error())
	} else {
		session.SetLastSuccessMsg(w, "You are now subscribed!")
	}
	http.Redirect(w, r, redirect, http.StatusTemporaryRedirect)

}

func HandleUnsubscripe(w http.ResponseWriter, r *http.Request) {

	conf, err := config.Config()
	if err != nil {
		log.Printf("Error loading config %s", err)
	}

	db, _ := cloudsql.CloudSQLConnection()
	defer db.Close()

	redirect := conf.StripeCardRedirect //os.Getenv("STRIPE_CARD_REDIRECT")
	if session.LastReferrerUrl(r) != "" {
		redirect = session.LastReferrerUrl(r)
		redirect, _ = url.QueryUnescape(redirect)
	}

	stripeMgr, err := NewStripeMgr(db)
	if err != nil {
		errS := "Error creating stripemgr: %s"
		log.Printf(errS, err)
		session.SetLastFailMsg(w,"Unable to remove subscription: " + err.Error())
		http.Redirect(w, r, redirect, http.StatusTemporaryRedirect)
	}

	prof := session.Profile(r)
	email := prof.Email
	if email == "" {
		panic("customer is not signed in")
	}
	subId := r.URL.Query().Get("subId")
	if subId == "" {
		panic("subId is required")
	}
	err = stripeMgr.UnsubscribeCustomer(subId)
	if err != nil {
		session.SetLastFailMsg(w,"Unable to remove subscription: " + err.Error())
	} else {
		session.SetLastSuccessMsg(w,"Successfully removed subscription.")
	}
	http.Redirect(w, r, redirect, http.StatusTemporaryRedirect)
}