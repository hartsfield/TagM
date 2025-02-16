package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

///////////////////////////////////////////////////////////////////////////////
////////////////////////      Password Section      ///////////////////////////
///////////////////////////////////////////////////////////////////////////////

// hashPassword takes a password string and returns a hash
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// checkPasswordHash compares a password to a hash and returns true if they
// match
func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

///////////////////////////////////////////////////////////////////////////////
////////////////////////       Auth Middleware      ///////////////////////////
///////////////////////////////////////////////////////////////////////////////
//     Every request made to an authenticated route is funneled through      //
//     this function:                                                        //
///////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////

// checkAuth parses and renews the authentication token, and adds it to the
// context. checkAuth is used as a middleware function for routes that allow or
// require authentication.
func checkAuth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get the "token" cookie
		token, err := r.Cookie("token")
		if err != nil {
			serveUnauthed(next, r, w, err)
			return
		}

		// parse the "token" cookie, making sure it's valid, and
		// obtaining user credentials if it is
		c, err := parseToken(token.Value)
		if err != nil {
			serveUnauthed(next, r, w, err)
			return
		}

		// check if "token" cookie matches the token stored in the
		// database
		if err = scanProfile(c); err != nil {
			serveUnauthed(next, r, w, err)
			return
		}
		// if the tokens match we renew the token and mark the user as
		// logged in
		if c.User.Token == token.Value {
			c.IsLoggedIn = true
			verified, err := renewToken(w, r, c)
			if err != nil {
				serveUnauthed(next, r, w, err)
				return
			}
			log.Println(c.User.ID)
			next.ServeHTTP(w, r.WithContext(verified))
			return
		}
		serveUnauthed(next, r, w, err)
	})
}

func serveUnauthed(next http.HandlerFunc, r *http.Request, w http.ResponseWriter, err error) {
	// create a generic user object thats not signed in to be used
	// as a placeholder until credentials are verified. Here, we
	// place it gently into the context using ctxkey (iota) as the
	// key corresponding to to the user value:
	if err != nil {
		log.Println(err)
	}
	unverified := context.WithValue(
		r.Context(),
		ctxkey,
		&credentials{IsLoggedIn: false},
	)
	next.ServeHTTP(w, r.WithContext(unverified))
}

///////////////////////////////////////////////////////////////////////////////
////////////////////////    JSON Web Token Stuff    ///////////////////////////
///////////////////////////////////////////////////////////////////////////////

// parseToken takes a token string, checks its validity, and parses it into a
// set of credentials. If the token is invalid it returns an error
func parseToken(tokenString string) (*credentials, error) {
	var claims *credentials = &credentials{IsLoggedIn: false, User: &user{Token: ""}}
	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			return hmacSampleSecret, nil
		},
	)
	if err == nil {
		if claims, ok := token.Claims.(*credentials); ok && token.Valid {
			return claims, nil
		}
	}
	return nil, err
}

// renewToken renews a users token using existing claims, sets it as a cookie
// on the client, and adds it to the database.
// TODO: FIX EXPIRY
func renewToken(w http.ResponseWriter, r *http.Request, c *credentials) (context.Context, error) {
	c.User.Token = ""
	ss, err := jwt.NewWithClaims(jwt.SigningMethodHS256,
		&credentials{c.Name, "", true, c.User,
			jwt.StandardClaims{
				ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
			},
		}).SignedString(hmacSampleSecret)
	if err != nil {
		return nil, err
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   ss,
		Path:    "/",
		Expires: time.Now().Add(10 * time.Minute),
		MaxAge:  0,
	})
	c.User.Token = ss
	if err = setProfile(c); err != nil {
		return nil, err
	}
	return context.WithValue(r.Context(), ctxkey, c), nil
}

///////////////////////////////////////////////////////////////////////////////
////////////////////////     Marshal Credentials    ///////////////////////////
///////////////////////////////////////////////////////////////////////////////

// marshalCredentials is used convert a request body into a credentials{}
// struct
func marshalCredentials(r *http.Request) (*credentials, error) {
	t := new(credentials)
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	err := decoder.Decode(t)
	if err != nil {
		return t, err
	}
	return t, nil
}
