package main

import (
	"log"
	"net/http"
	"regexp"
)

func status(w http.ResponseWriter, s string, err error) error {
	ajaxResponse(w, map[string]string{
		"status": s,
	})
	return err
}

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
