// Package auth0 provides a set of handlers for integrating Auth0's (https://auth0.com/) service with Go Web Applications running in Google App Engine.
package auth0

import (
	_ "crypto/sha512"
	"encoding/json"
	"golang.org/x/oauth2"
	"io/ioutil"
	"net/http"
	"github.com/ezeev/saga/session"
	cnprofile "github.com/ezeev/saga/profile"
	"github.com/ezeev/saga/util"
	"net/url"
	"github.com/ezeev/saga/metrics"
	"github.com/ezeev/saga/config"
)


var reqVars = []string{"AUTH0_CALLBACK_URI","AUTH0_SIGNOUT_URI","AUTH0_CALLBACK_HOST_DEV","AUTH0_CALLBACK_HOST_LIVE","AUTH0_CLIENT_ID","AUTH0_CLIENT_SECRET"}

// RegisterHandlers Is a helper function that will register all Auth0 handlers using the options set in
// your app.yaml
func init() {

	err := util.CheckVars(reqVars)
	if err != nil {
		metrics.Registry().IncAuth0Errors()
		panic(err)
	}
	sconf, err := config.Config()
	if err != nil {
		panic(err)
	}
	http.HandleFunc(sconf.Auth0CallbackURI,CallbackHandler)
	http.HandleFunc(sconf.Auth0SignoutURI,HandleSignout)
}

// CallbackHandler handles the Auth0 callback. After completing the Auth0 handshake it
// will start the user's session and redirect back to the original referrer or the path set
// in the OAUTH_SUCCESS_REDIRECT environment variable.
func CallbackHandler(w http.ResponseWriter, r *http.Request) {


	//someone is attempting to login
	metrics.Registry().IncLoginAttempts()
	sconf, _ := config.Config()

	domain := sconf.Auth0Domain

	var callBackUrl string
	if sconf.IsDev {
		callBackUrl = sconf.Auth0CallbackHostDev + sconf.Auth0CallbackURI
	} else {
		callBackUrl = sconf.Auth0CallbackHostLive + sconf.Auth0CallbackURI
	}


	conf := &oauth2.Config{
		ClientID:     sconf.Auth0ClientID,
		ClientSecret: sconf.Auth0ClientSecret,
		RedirectURL:  callBackUrl,
		Scopes:       []string{"openid", "profile"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://" + domain + "/authorize",
			TokenURL: "https://" + domain + "/oauth/token",
		},
	}

	code := r.URL.Query().Get("code")
	token, err := conf.Exchange(oauth2.NoContext, code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Getting now the userInfo
	client := conf.Client(oauth2.NoContext, token)
	resp, err := client.Get("https://" + domain + "/userinfo")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	raw, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var profile map[string]interface{}
	if err = json.Unmarshal(raw, &profile); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var email, photo string
	if str, ok := profile["email"].(string); ok {
		email = str
	}
	if str, ok := profile["picture"].(string); ok {
		photo = str
	}

	if email == "" {
		session.SetLastFailMsg(w,"You cannot use this account to authenticate. We require an email address. Either there is no email address associated with this account or the account is not verified or your email is not accessible to auth providers.")
		http.Redirect(w, r, sconf.OAuthSuccessRedirect, http.StatusTemporaryRedirect)
	}
	//create a new session
	prof := cnprofile.NewProfile(email,photo,token.AccessToken)
	//session.Start(w,email,photo,token.AccessToken)
	session.Start(w,prof)

	// decide where to redirect the user
	redirect := sconf.OAuthSuccessRedirect
	if session.LoginReferrerUrl(r) != "" {
		redirect = session.LoginReferrerUrl(r)
		redirect, _ = url.QueryUnescape(redirect)
	} else if session.LastReferrerUrl(r) != "" {
		redirect = session.LastReferrerUrl(r)
		redirect, _ = url.QueryUnescape(redirect)
	}
	// Redirect to logged in page
	http.Redirect(w, r, redirect, http.StatusSeeOther)
}


// HandleSignout expires the user's session and redirects to the referrer
// or to the path specified in the OAUTH_SIGNOUT_REDIRECT environment variable.
func HandleSignout(w http.ResponseWriter, r *http.Request) {
	session.End(w)

	sconf, _ := config.Config()

	redirect := sconf.Auth0SignoutRedirectURI
	if session.LastReferrerUrl(r) != "" {
		redirect = session.LastReferrerUrl(r)
		redirect, _ = url.QueryUnescape(redirect)
	}
	http.Redirect(w,r,redirect, http.StatusTemporaryRedirect)

}