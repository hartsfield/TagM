package main

import (
	"log"
	"net/http"
	"time"
)

// signout logs the user out by overwriting the token. It must first validate
// the existing token to get the username to overwrite the old token in the
// database
func signout(w http.ResponseWriter, r *http.Request) {
	log.Println(" testfpoemvpofremvop")
	// token, err := r.Cookie("token")
	// if err != nil {
	// 	log.Println(err)
	// }
	//
	// _, err = parseToken(token.Value)
	// if err != nil {
	// 	log.Println(err)
	// }
	// c.User.Token = ""
	// setProfile(c)
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   "loggedout",
		Path:    "/",
		Expires: time.Now(),
		MaxAge:  0,
	})

	w.Write(nil)
}
