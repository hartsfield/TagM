package main

import (
	"log"
	"net/http"
)

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
