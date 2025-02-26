// / Provided Under BSD (2 Clause)
//
// Copyright 2025 Johnathan A. Hartsfield
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice,
// this list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice,
// this list of conditions and the following disclaimer in the documentation
// and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS “AS IS”
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
// ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
// LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
// CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
// SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
// INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
// CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
// ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
// POSSIBILITY OF SUCH DAMAGE.
//
// ////////////////////////////////////////////////////////////////////////////
//
// auth_handlers.go houses identification, authentication, and token
// re-newel functionality. We use JSON Web Tokens (JWT) stored as a cookie in
// the clients http request header, and in the database for re/authentication.
// Passwords are never stored in plaintext, and are instead stored in hashed
// form using bcrypt.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

///////////////////////////////////////////////////////////////////////////////
////////////////////////      Password Section      ///////////////////////////
///////////////////////////////////////////////////////////////////////////////

// hashPassword() takes a password string and returns a hash
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// checkPasswordHash() compares a password to a hash and returns true if they
// match
func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

///////////////////////////////////////////////////////////////////////////////
////////////////////////       Signup Section       ///////////////////////////
///////////////////////////////////////////////////////////////////////////////

// signup signs a user up. It's a response to an XMLHttpRequest (AJAX request)
// containing new user credentials. It responds with a map[string]string that
// can be converted to JSON.
func signup(w http.ResponseWriter, r *http.Request) {
	// Marshal the Credentials into a credentials struct
	c, err := marshalCredentials(r)
	if err != nil {
		log.Println(status(w, "Invalid Credentials", err))
		return
	}

	// Make sure the username doesn't contain forbidden symbols
	emailRegx := "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"
	match, err := regexp.MatchString(emailRegx, c.Name)
	if err != nil || !match {
		log.Println(status(w, "Invalid Username (E1)", err))
		return
	}

	// TODO: Move to dbcalls.go
	if rdb.Exists(rdx, c.Name+":HASH").Val() > 0 {
		log.Println(status(w, "User Exists", nil))
		return
	}

	// Check to make sure the password is long enough. We don't store
	// sensitive info or recommend users upload it, so seven works.
	if len(c.Password) < 7 {
		log.Println(status(w, "Invalid Password (E2)", nil))
		return
	}

	// Save an ID to be used for posting so we don't expose
	// the users email, save this in redis as an HSET, and
	// store all the profile information here.
	if c.User == nil {
		c.User = new(user)
	}
	c.User.ID = genID(15)
	c.User.ProfileBG = "public/media/hubble.jpg"
	c.User.ProfilePic = "public/media/ndt.jpg"

	// If username is unique and valid, we attempt to hash
	// the password
	hash, err := hashPassword(c.Password)
	if err != nil {
		log.Println(status(w, "Invalid Password", err))
		return
	}
	c.Password = ""
	_, err = setHashToID(c, hash)
	if err != nil {
		log.Println(status(w, "Database Error", err))
		return
	}

	// If the password is hashable, and we were able to add
	// the user to the redis ZSET, we store the hash in the
	// database with the username as the key and the hash
	// as the value thats returned by the key.
	if _, err = setPasswordHash(c, hash); err != nil {
		log.Println(status(w, "Database Error", err))
		return
	}

	// Add the user the USERS set in redis. This
	// associates a score with the user that can be
	// incremented or decremented
	if _, err = zaddUsers(c); err != nil {
		log.Println(status(w, "Database Error", err))
		return
	}

	// renewToken() is used here to create a new token, which it's also
	// capable of.
	if _, err = renewToken(w, r, c); err != nil {
		log.Println(status(w, "Token Error", err))
		return
	}

	// success.
	log.Println(status(w, "success", nil))
}

///////////////////////////////////////////////////////////////////////////////
////////////////////////       Signin Section       ///////////////////////////
///////////////////////////////////////////////////////////////////////////////

// signin signs a user in. It's a response to an XMLHttpRequest (AJAX request)
// containing the user credentials. It responds with a map[string]string that
// can be converted to JSON by the client. The client expects a boolean
// indicating success or error, and a possible error string.
func signin(w http.ResponseWriter, r *http.Request) {
	// Marshal the Credentials into a credentials struct
	c, err := marshalCredentials(r)
	if err != nil {
		log.Println(status(w, "Invalid Credentials", err))
		return
	}

	// Get the passwords hash from the database by looking up the users
	// name
	hash, err := getPasswordHash(c)
	if err != nil {
		log.Println(status(w, "User doesn't exist", err))
		return
	}

	// Check if password matches by hashing it and comparing the hashes
	doesMatch := checkPasswordHash(c.Password, hash)
	if doesMatch {
		c.Password = "" // remove the password from credentials{}
		// before doing anything else.

		// Get the users ID from the database using the hash obtained
		// previously.
		id, err := getID(hash)
		if err != nil {
			log.Println(status(w, "Database Error", err))
			return
		}

		// Create a user{} with an ID set to the ID obtained
		// previously.
		c.User = &user{ID: id}

		// Use the previously created user{} to get the rest of the
		// users profile information.
		err = scanProfile(c)
		if err != nil {
			log.Println(status(w, "Scan Profile Error", err))
			return
		}

		// issue/renew the users authentication token.
		_, err = renewToken(w, r, c)
		if err != nil {
			log.Println(status(w, "Token Error", err))
			return
		}

		// success
		log.Println(status(w, "success", err))
		return
	}
	log.Println(status(w, "Bad Password", err))
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
		log.Println(token)
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
			// success
			next.ServeHTTP(w, r.WithContext(verified))
			return
		}
		serveUnauthed(next, r, w, err)
	})
}

// serveUnauthed() is used when a user fails an authentication challenge and so
// is served with a page a user without an account would see.
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
	// Create a dummy credentials{}. We implemented jwt.StandardClaims on
	// the credentials{} struct defined in main.go so that we may pass it
	// to jwt.ParseWithClaims(). This is for convenience sake.
	var claims *credentials = &credentials{
		IsLoggedIn: false,
		User:       &user{Token: ""},
	}

	// Use the json web token module function jwt.ParseWithClaims() to
	// parse the token passed herein.
	token, err := jwt.ParseWithClaims(tokenString, claims, jwtkey_fn)
	if err != nil {
		return nil, err
	}
	var ok bool
	if claims, ok = token.Claims.(*credentials); !ok || !token.Valid {
		return nil, errors.New("Invalid Token or Type Assertion")
	}

	// success
	return claims, nil
}

func jwtkey_fn(token *jwt.Token) (interface{}, error) {
	// hmacSampleSecret should be defined at run time as an
	// environment variable. see: main.go.
	return hmacSampleSecret, nil
}

// renewToken renews a users token using existing claims, sets it as a cookie
// on the client, and adds it to the database.
// TODO: FIX EXPIRY
func renewToken(w http.ResponseWriter, r *http.Request, c *credentials) (context.Context, error) {
	c.User.Token = "" // make sure the old token is removed.

	// use the functionality provided by the json web token module to renew
	// the token using the jwt.StandardClaims{}.
	ss, err := jwt.NewWithClaims(jwt.SigningMethodHS256,
		&credentials{c.Name, "", true, c.User,
			jwt.StandardClaims{
				ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
			},
		}).SignedString(hmacSampleSecret)
	if err != nil {
		return nil, err
	}

	// Set the token as a cookie in the response headers.
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   ss,
		Path:    "/",
		Expires: time.Now().Add(10 * time.Minute),
		MaxAge:  0,
	})

	// Update the users token in the database.
	c.User.Token = ss
	if err = setProfile(c); err != nil {
		return nil, err
	}

	// success
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
