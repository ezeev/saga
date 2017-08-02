// Package ratelimit provides a method for incrementing an IP Address based rate limiter
package ratelimit

import (
	"google.golang.org/appengine/memcache"
	"net/http"
	"time"
	"google.golang.org/appengine"
	"strconv"
	"fmt"
)

const cachePrefix = "RATE-"


// Increment increments an IP Address based rate limiter metric stored in Google Cloud's Memcache
func Increment(r *http.Request, limit uint64) (string, uint64, error) {
	c := appengine.NewContext(r)
	unixIntValue := time.Now().Unix()
	//fmt.Println(timeStamp)
	timeStamp := time.Unix(unixIntValue, 0)
	hr, min, _ := timeStamp.Clock()
	key := r.RemoteAddr+"-"+strconv.Itoa(hr)+"-"+strconv.Itoa(min)
	val, err := memcache.Increment(c,key,1,0)
	if val >= limit {
		return "",val, &RateLimitError{msg: fmt.Sprintf("IP Address %s has exceeded the rate limit of %d requests per minute",r.RemoteAddr,limit)}
	}
	return key, val, err
}

type RateLimitError struct {
	msg string

}

func (e *RateLimitError) Error() string { return e.msg }