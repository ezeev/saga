// Package profile provides a struct representing a user's profile and methods for serializing
// and de-serializing a profile as Java Web Token (jwt).
package profile

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
	"encoding/json"
)

type Profile struct {
	Email string
	Photo string
	AccessToken string
}

var signingString = []byte("thisismysecret")

// NewProfile creates a new profile and returns a pointer to it.
func NewProfile(email string, photo string, accessToken string) *Profile {
	profile := &Profile{}
	profile.Email = email
	profile.Photo = photo
	profile.AccessToken = accessToken
	return profile
}

// Serializes a profile struct to a Java Web Token.
func ToJwt(profile *Profile) (string, error) {

	profData, err := json.Marshal(&profile)
	if (err != nil) {
		return "", fmt.Errorf("Unable to marshall profile to json")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"profile": string(profData),
		"created": time.Now().Unix(),
	})
	tokenString, err := token.SignedString(signingString)
	if err != nil {
		fmt.Println("Error in GenerateJwtToken")
		fmt.Print(err)
	}
	return tokenString, nil
}

// De-serializes a Java Web Token to a profile struct and an epoch timestamp of when the token was created.
func ToProfile(tokenString string) (*Profile, int64, error) {
	var resp Profile
	var created int64

	if tokenString == "" {
		return nil,0,nil
	}
	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	parser := jwt.Parser{}
	parser.UseJSONNumber = true
	//token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
	token, err := parser.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return signingString, nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if str, ok := claims["profile"].(string); ok {
			profileData := []byte(str)
			err := json.Unmarshal(profileData, &resp)
			if err != nil {
				return nil,0,fmt.Errorf("Unable to parse profile data")
			}
		} else {
			fmt.Println("profile claim is not a string")
		}
		if ts, ok := claims["created"].(json.Number); ok {
			created, err = ts.Int64()
			if err != nil {
				fmt.Println("Unable to parse created ts from jwt claim")
			}
		} else {
			fmt.Println("created timestamp claim is not an int64")
		}


		return &resp, created, nil
	} else {
		return nil, 0, err
	}
	return &resp, created, nil
}