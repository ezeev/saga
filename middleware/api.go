// Package middleware provides http HandlerFuncs for use with http web or API endpoints.
package middleware

import (
	"fmt"
	"net/http"
	"github.com/ezeev/saga/profile"
	"github.com/ezeev/saga/ratelimit"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"github.com/ezeev/saga/config"
)

// ApiAuth, when used on an API endpoint, provides CORS (including "OPTIONS" request support),
// Java Web Token (jwt) validation/authorization, and IP Address based rate limiting.
func ApiAuth(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		conf, err := config.Config()
		if err != nil {
			panic(err)
		}

		//Set CORS Headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, X-Auth-Token")
		if r.Method == "OPTIONS" {
			//don't need to do anything else
			return
		}
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
			// rate limit
			//rateLimit := os.Getenv("API_RATE_LIMIT")
			//limiter := time.Tick(time.Millisecond * 200)
			c := appengine.NewContext(r)
			intlimit := conf.ApiRateLimitPerMin
			key, count, err := ratelimit.Increment(r,uint64(intlimit)	)
			if err != nil {
				w.WriteHeader(http.StatusTooManyRequests)
				fmt.Fprintln(w, err.Error())
				return
			}
			log.Infof(c,"Hit Counter:%s %d", key, count)
			fn.ServeHTTP(w, r)
		}
	}
}
