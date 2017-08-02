
// Package session provides session and cookie management. This package works with other cloud-ninja packages
// but can easily be used independently with other packages.
package session

import (
	"net/http"
	"time"
	"github.com/ezeev/saga/profile"
)

const loginRefCookieId = "cn-loginref"
const lastSuccessMsgCookieId = "cn-lastsuccessmsg"
const lastFailMsgCookieId = "cn-lastfailmsg"
const profileCookieId = "cn-profile"
const emailCookieId = "cn-email"
const lastReferrerCookieId = "cn-ref"

func SetLoginReferrerUrl(w http.ResponseWriter, path string) {
	saveCookie(w,loginRefCookieId,path)
}

func SetLastReferrerUrl(w http.ResponseWriter, path string) {
	saveCookie(w,lastReferrerCookieId,path)
}

func SetLastSuccessMsg(w http.ResponseWriter, msg string) {
	saveCookie(w,lastSuccessMsgCookieId,msg)
}

func LastSuccessMsg(w http.ResponseWriter, r *http.Request) string {
	msg := getCookie(r,lastSuccessMsgCookieId)
	expCookie(w,lastSuccessMsgCookieId)
	return msg
}


func SetLastFailMsg(w http.ResponseWriter, msg string) {
	saveCookie(w,lastFailMsgCookieId,msg)
}

func LastFailMsg(w http.ResponseWriter, r *http.Request) string {
	msg := getCookie(r,lastFailMsgCookieId)
	expCookie(w,lastFailMsgCookieId)
	return msg
}

// Start will start a new user session. It uses a profile struct as input.
func Start(w http.ResponseWriter, prof *profile.Profile) {

	jwt, err := profile.ToJwt(prof)
	if err != nil {
		panic("Unable to create Jwt token from profile struct")
	}
	saveCookie(w,profileCookieId,jwt)
	saveCookie(w,emailCookieId,prof.Email)

}

// End ends a user's session by expiring their cookies
func End(w http.ResponseWriter) {
	expCookie(w,profileCookieId)
	expCookie(w,loginRefCookieId)
	expCookie(w,emailCookieId)
}

// Profile retrieves a user's JWT token cookie and then converts it to a profile struct.
func Profile(r *http.Request) *profile.Profile {
	jwt := getCookie(r,profileCookieId)
	prof,_,err := profile.ToProfile(jwt)

	if err != nil {
		panic("Unable to parse jwt token into profile struct")
	}
	return prof
}


func UpdateProfile(w http.ResponseWriter, r *http.Request, prof *profile.Profile) {
	cookie, _ := r.Cookie(profileCookieId)
	cookie.Value, _ = profile.ToJwt(prof)
	http.SetCookie(w,cookie)
}


func LoginReferrerUrl(r *http.Request) string {
	return getCookie(r, loginRefCookieId)
}

func LastReferrerUrl(r *http.Request) string {
	return getCookie(r, lastReferrerCookieId)
}

func saveCookie(w http.ResponseWriter, cookieId string, value string) {
	expiration := time.Now().Add(time.Hour * 336) // 2 week exp
	cookie := http.Cookie{Name: cookieId, Value:value, Expires: expiration}
	http.SetCookie(w, &cookie)
}

func getCookie(r *http.Request, cookieId string) string {
	cookie, _ := r.Cookie(cookieId)
	if cookie != nil {
		return cookie.Value
	} else {
		return ""
	}
}

func expCookie(w http.ResponseWriter, cookieId string) {
	expiration := time.Now().Truncate(time.Hour)
	cookie := http.Cookie{Name: cookieId, Value:"", Expires: expiration}
	http.SetCookie(w, &cookie)
}