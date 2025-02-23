package main

import (
	"context"
	"encoding/json"
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
	log.Println(c)
	// Make sure the username doesn't contain forbidden symbols
	emailRegx := "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"
	match, err := regexp.MatchString(emailRegx, c.Name)
	if err != nil || !match {
		log.Println(status(w, "Invalid Username (E1)", err))
		return
	}
	if rdb.Exists(rdx, c.Name+":HASH").Val() > 0 {
		log.Println(status(w, "User Exists", nil))
		return
	}

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
		log.Println(err)
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

	if _, err = renewToken(w, r, c); err != nil {
		log.Println(status(w, "Token Error", err))
	}
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
		c.Password = ""
		id, err := getID(hash)
		if err != nil {
			log.Println(err)
		}
		c.User = &user{ID: id}

		err = scanProfile(c)
		if err != nil {
			log.Println(status(w, "Scan Profile Error", err))
			return
		}
		_, err = renewToken(w, r, c)
		if err != nil {
			log.Println(err)
			return
		}
		log.Println(status(w, "success", err))
		return
	}
	ajaxResponse(w, map[string]string{"success": "false", "error": "Bad Password"})
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
