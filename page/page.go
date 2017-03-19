// Package page aims to provide a single struct with all information needed by web applications about the user's state
// and data needed to render a web page or API response.
package page

import (
	"github.com/stripe/stripe-go"
	"github.com/ezeev/saga/session"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	stripeManager "github.com/ezeev/saga/stripe"
	"github.com/ezeev/saga/profile"
	"net/http"
	"os"
	"fmt"
	"html/template"
	"time"
	"database/sql"
)

// PageFuncMap provides a set of functions that can be passed to Go templates for displaying content
// to end users.
func PageFuncMap() template.FuncMap {
	fm := template.FuncMap{
		"formatTime": func(ts int64) string {
			tm := time.Unix(ts, 0)
			return tm.Format(time.RFC3339)
		},
		"displayAmount" : func(amount uint64) string {
			str := fmt.Sprint(amount)
			str = "$" + str[:len(str) - 2] + "." + str[len(str) - 2:]
			return str
		},
		"multiply" : func(x uint64, y uint64) uint64 {
			return x * y
		},
	}
	return fm
}


// Page is a struct composed of all data needed to render pages or API responses for a user.
type Page struct {
	Title string
	UserProfile *profile.Profile
	StripeId string
	StripeCustomer *stripe.Customer
	Plans *[]stripe.Plan
	Cards *[]stripe.Card
	LastSuccessMsg string
	LastFailMsg string
	StripePubKey string
	Auth0CallBackUrl string
	Path string
}

// NewPage creates a new Page struct and returns a pointer to it.
// If the user is logged in, it will load the respective child structs for the user's account.
func NewPage(w http.ResponseWriter, r *http.Request, db *sql.DB, jwt string) (*Page, error) {

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
			panic(err)
		}
	}
	page.UserProfile = prof

	log.Infof(c,"Profile: %s",prof)
	// make sure they have a stripe Id if logged in

	if page.UserProfile != nil {
		stripeMgr, err := stripeManager.NewStripeMgr(c,db)
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
			planIds = append(planIds,v.Plan.ID)
		}

	}

	page.LastFailMsg = session.LastFailMsg(w,r)
	page.LastSuccessMsg = session.LastSuccessMsg(w,r)

	var pubKey string
	if appengine.IsDevAppServer() {
		pubKey = os.Getenv("STRIPE_TEST_PK")
	} else {
		pubKey = os.Getenv("STRIPE_LIVE_PK")
	}

	page.StripePubKey = pubKey

	var callBackUrl string

	if appengine.IsDevAppServer() {
		callBackUrl = os.Getenv("AUTH0_CALLBACK_HOST_DEV") + os.Getenv("AUTH0_CALLBACK_URI")
	} else {
		callBackUrl = os.Getenv("AUTH0_CALLBACK_HOST_LIVE") + os.Getenv("AUTH0_CALLBACK_URI")
	}
	page.Auth0CallBackUrl = callBackUrl

	return page, nil
}
