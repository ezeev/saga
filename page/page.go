// Package page aims to provide a single struct with all information needed by web applications about the user's state
// and data needed to render a web page or API response.
package page

import (
	"database/sql"
	"fmt"
	"github.com/ezeev/saga/metrics"
	"github.com/ezeev/saga/profile"
	"github.com/ezeev/saga/session"
	stripeManager "github.com/ezeev/saga/stripe"
	"github.com/stripe/stripe-go"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"html/template"
	"net/http"
	"time"
	"github.com/ezeev/saga/config"
)

// PageFuncMap provides a set of functions that can be passed to Go templates for displaying content
// to end users.
func PageFuncMap() template.FuncMap {
	fm := template.FuncMap{
		"formatTime": func(ts int64) string {
			tm := time.Unix(ts, 0)
			return tm.Format(time.RFC3339)
		},
		"displayAmount": func(amount uint64) string {
			str := fmt.Sprint(amount)
			str = "$" + str[:len(str)-2] + "." + str[len(str)-2:]
			return str
		},
		"multiply": func(x uint64, y uint64) uint64 {
			return x * y
		},
	}
	return fm
}

// Page is a struct composed of all data needed to render pages or API responses for a user.
type Page struct {
	Title             string
	UserProfile       *profile.Profile
	StripeId          string
	StripeCustomer    *stripe.Customer
	Plans             *[]stripe.Plan
	Cards             *[]stripe.Card
	LastSuccessMsg    string
	LastFailMsg       string
	StripePubKey      string
	Auth0CallBackUrl  string
	Path              string
	Auth0ClientId     string
	Auth0Domain       string
	Auth0CallBackHost string
	Auth0CallBackURI  string
	AppDomain         string
	AppName           string
}

// NewPage creates a new Page struct and returns a pointer to it.
// If the user is logged in, it will load the respective child structs for the user's account.
func NewPage(w http.ResponseWriter, r *http.Request, db *sql.DB, jwt string) (*Page, error) {

	conf, _ := config.Config()

	c := appengine.NewContext(r)
	page := &Page{}
	page.Path = r.RequestURI
	var prof *profile.Profile
	var err error
	if jwt == "" {
		prof = session.Profile(r)
	} else {
		prof, _, err = profile.ToProfile(jwt)
		if err != nil {
			log.Errorf(c, "Error in NewPage while getting session: %s", err)
			metrics.Registry().IncPageLoadErrors()
			panic(err)
		}
	}
	page.UserProfile = prof
	page.AppName = conf.AppName
	page.Auth0ClientId = conf.Auth0ClientID
	page.Auth0CallBackURI = conf.Auth0CallbackURI
	page.AppDomain = conf.AppDomain
	if appengine.IsDevAppServer() {
		page.Auth0CallBackHost = conf.Auth0CallbackHostDev
	} else {
		page.Auth0CallBackHost = conf.Auth0CallbackHostLive
	}

	if page.UserProfile != nil {
		stripeMgr, err := stripeManager.NewStripeMgr(c, db)
		if err != nil {
			return nil, err
		}
		email := page.UserProfile.Email
		cust, err := stripeMgr.Customer(email)
		if err != nil {
			return nil, err
		}
		page.StripeCustomer = cust
		page.StripeId = cust.ID
		page.Plans, err = stripeMgr.ListPlans()
		if err != nil {
			return nil, err
		}
		//populate the customer's cards
		page.Cards, err = stripeMgr.ListCards(cust.ID)
		if err != nil {
			return nil, err
		}

		//store the user's subscribed plans in their jwt token
		var planIds []string
		for _, v := range cust.Subs.Values {
			planIds = append(planIds, v.Plan.ID)
		}

	}

	page.LastFailMsg = session.LastFailMsg(w, r)
	page.LastSuccessMsg = session.LastSuccessMsg(w, r)

	var pubKey string
	if appengine.IsDevAppServer() {
		pubKey = conf.StripeTestPublicKey
	} else {
		pubKey = conf.StripeLivePublicKey
	}

	page.StripePubKey = pubKey

	var callBackUrl string

	if appengine.IsDevAppServer() {
		callBackUrl = conf.Auth0CallbackHostDev + conf.Auth0CallbackURI
	} else {
		callBackUrl = conf.Auth0CallbackHostLive + conf.Auth0CallbackURI
	}

	page.Auth0Domain = conf.Auth0Domain

	page.Auth0CallBackUrl = callBackUrl
	return page, nil
}
