// Package stripe provides API endpoints for common Stripe operations.
// It provides wrapper methods around the Stripe API. It can be used with
// other cloud-ninja packages but there are no direct dependencies to other
// cloud-ninja packages.
package stripe

import (
	"github.com/stripe/stripe-go/client"
	"github.com/stripe/stripe-go"
	//"google.golang.org/appengine/urlfetch"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"os"
	"golang.org/x/net/context"
	"net/http"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"crypto/tls"
)


type StripeMgr struct {
	PubApiKey string
	SecApiKey string
	AppFilter string
	AppEngineContext context.Context
	httpClient *http.Client
	stripeClient *client.API
	db *sql.DB
}



func NewStripeMgr(appEngineContext context.Context, db *sql.DB) (*StripeMgr, error) {
	if os.Getenv("STRIPE_TEST_SK") == "" || os.Getenv("STRIPE_TEST_PK") == "" {
		log.Errorf(appEngineContext,"STRIPE_TEST_SK and STRIPE_TEST_PK env vars must be set!")
		return nil, fmt.Errorf("STRIPE_TEST_SK and STRIPE_TEST_PK env vars must be set!")
	}
	var pubKey string
	var secKey string

	if appengine.IsDevAppServer() {
		pubKey = os.Getenv("STRIPE_TEST_PK")
		secKey = os.Getenv("STRIPE_TEST_SK")
	} else {
		pubKey = os.Getenv("STRIPE_LIVE_PK")
		secKey = os.Getenv("STRIPE_LIVE_SK")
	}

	if pubKey == "" || secKey == "" {
		return nil, fmt.Errorf("Stripe Keys are net set, please set STRIPE_LIVE_PK and STRIPE_LIVE_SK or replace LIVE with TEST if in test mode.")
	}

	mgr := &StripeMgr{
		PubApiKey: pubKey,
		SecApiKey: secKey,
		AppEngineContext: appEngineContext,
	}
	mgr.AppFilter = os.Getenv("STRIPE_APP_FILTER")

	//Stripe now requires TLS 1.2
	//mgr.httpClient = urlfetch.Client(mgr.AppEngineContext)
	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{},
		DisableCompression: true,
	}

	mgr.httpClient = &http.Client{Transport: tr}
	mgr.stripeClient = client.New(mgr.SecApiKey,stripe.NewBackends(mgr.httpClient))
	mgr.db = db
	return mgr, nil
}


func (this *StripeMgr) getStripeIdDB(email string) (string, error) {

	var env string
	if appengine.IsDevAppServer() {
		env = "dev"
	} else {
		env = "live"
	}

	var stripeId string
	err := this.db.QueryRow("select stripe_id from customers where email = ? and env = ?", email, env).Scan(&stripeId)
	if err != nil {
		this.logError(err)
		return "", err
	}
	return stripeId, nil
}

func (this *StripeMgr) saveStripeIdDB(email string, stripeId string) error {

	var env string
	if appengine.IsDevAppServer() {
		env = "dev"
	} else {
		env = "live"
	}

	stmt, err := this.db.Prepare("INSERT INTO customers(email,stripe_id,env) VALUES(?,?,?)")
	if err != nil {
		this.logError(err)
		return err
	}
	_, err = stmt.Exec(email,stripeId,env)
	if err != nil {
		this.logError(err)
		return err
	}
	return nil
}

func (this *StripeMgr) IsSubscribed(email string, planId string) (bool, error) {
	cust, err := this.Customer(email)
	if err != nil {
		return false, err
	}

	log.Infof(this.AppEngineContext,"Here's all of the subs for %s: %s", email, cust.Subs.Values)

	for _, v := range cust.Subs.Values {
		if v.Plan.ID == planId {
			return true, nil
		}
	}
	return false, nil
}

func (this *StripeMgr) logError(err error) {
	log.Errorf(this.AppEngineContext, "Error in StripeMgr: %s", err)
}

func (this *StripeMgr) paramsFilter() map[string]string {
	return map[string]string{"app": this.AppFilter}
}


func (this *StripeMgr) createCustomer(email string) (*stripe.Customer, error) {

	customerParams := &stripe.CustomerParams{
		Desc: "New Customer",
		Email: email,
	}
	if this.AppFilter != "" {
		customerParams.Meta = this.paramsFilter()
	}

	c, err := this.stripeClient.Customers.New(customerParams)
	if err != nil {
		return nil, err

	}
	//stripeId successfully added, save it in db
	this.saveStripeIdDB(email,c.ID)

	return c, nil
}

func (this *StripeMgr) getCustomer(stripeId string) (*stripe.Customer, error) {
	//query stripe to make sure it is a valid customer
	c, err := this.stripeClient.Customers.Get(stripeId,nil)
	if err != nil {
		this.logError(err)
		return c, err
	}
	return c, nil
}

func (this *StripeMgr) ListPlans() (*[]stripe.Plan, error) {
	params := &stripe.PlanListParams{}
	//params.Filters.AddFilter("limit", "", "3")
	plans := this.stripeClient.Plans.List(params)
	var arrPlans []stripe.Plan

	for plans.Next() {
		//filter the plans by our app
		if this.AppFilter != "" {
			if plans.Plan().Meta["app"] == this.AppFilter {
				arrPlans = append(arrPlans, *plans.Plan())
			}
		} else {
			arrPlans = append(arrPlans, *plans.Plan())
		}
	}
	return &arrPlans, nil
}

func (this *StripeMgr) ListCards(custId string) (*[]stripe.Card, error) {
	params := &stripe.CardListParams{
		Customer: custId,
	}

	//params.Filters.AddFilter("limit", "", "3")
	cards := this.stripeClient.Cards.List(params)
	var arrCards []stripe.Card
	for cards.Next() {
		arrCards = append(arrCards,*cards.Card())
	}
	return &arrCards, nil
}

func (this *StripeMgr) AddCard(custId string, token string) (*stripe.Card, error) {

	params := stripe.CardParams{
		Customer: custId,
		Token: token,
	}
	c, err := this.stripeClient.Cards.New(&params)
	if err != nil {
		return nil, err
	}
	return c, nil

}

func (this *StripeMgr) DelCard(custId string, cardId string) (*stripe.Card, error) {
	c, err := this.stripeClient.Cards.Del(cardId, &stripe.CardParams{Customer: custId})
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (this *StripeMgr) Customer(email string) (*stripe.Customer, error) {

	//does the customer already exist in the database?
	stripeId, _ := this.getStripeIdDB(email)
	/*if err != nil {
		// don't throw error necessarily if no id exists
		return nil, nil
	}*/

	//yes, get the customer payload from stripe
	if stripeId != "" {
		c, err := this.getCustomer(stripeId)
		if err != nil {
			//don't necessarily throw an error here
			this.logError(err)
		}
		return c, nil
	} else {
		c, err := this.createCustomer(email)
		if err != nil {
			return nil, err
		}
		return c, nil
	}
}

func (this *StripeMgr) SubscribeCustomer(email string, planId string) (*stripe.Sub, error) {

	stripeId, err := this.getStripeIdDB(email)
	if err != nil {
		return nil, err
	}

	params := &stripe.SubParams{
		Customer: stripeId,
		Plan: planId,
	}

	s, err := this.stripeClient.Subs.New(params)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (this *StripeMgr) UnsubscribeCustomer(subId string) error {

	params := &stripe.SubParams{}
	_, err := this.stripeClient.Subs.Cancel(subId,params)
	if err != nil {
		return err
	}
	return nil
}