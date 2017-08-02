

// Package middleware provides http HandlerFuncs for use with http web or API endpoints.
package middleware

import (
	"fmt"
	"net/http"
	"github.com/ezeev/saga/profile"
	//"github.com/ezeev/saga/ratelimit"
)

// Api, when used on an API endpoint, provides CORS (including "OPTIONS" request support).
func Api(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//Set CORS Headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, X-Auth-Token")
		if r.Method == "OPTIONS" {
			//don't need to do anything else
			return
		}
		fn.ServeHTTP(w, r)
	}
}

// ApiAuth, when used on an API endpoint, provides Java Web Token (jwt) validation/authorization
func ApiAuth(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//look for the jwt header
		var jwtToken string
		jwtToken = r.Header.Get("X-Auth-Token")
		profile, _, err := profile.ToProfile(jwtToken)
		if err != nil {
			panic(err.Error())
		}
		if profile == nil {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintln(w, "{status: \"Unauthorized\"}")
		} else if profile.Email == "" {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintln(w, "{status: \"Unauthorized\"}")
		} else {
			fn.ServeHTTP(w, r)
		}
	}
}

/*func ApiRateLimit(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		conf, err := config.Config()
		if err != nil {
			panic(err)
		}

		// rate limiter
		c := appengine.NewContext(r)
		intlimit := conf.ApiRateLimitPerMin
		key, count, err := ratelimit.Increment(r,uint64(intlimit)	)
		if err != nil {
			w.WriteHeader(http.StatusTooManyRequests)
			fmt.Fprintln(w, err.Error())
			return
		}
		log.Debugf(c,"Hit Counter:%s %d", key, count)
		fn.ServeHTTP(w, r)
	}
}*/

