// helpers.go houses helper functions used either for initializations, or
// here and there throughout the tagmachine program. The functions in this file
// should be kept in order of importance.
package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"

	"github.com/redis/go-redis/v9"
)

// exeTmpl() is used to build and execute an html template, and is used in all
// the handlers which aren't ajax responders. This should be considered an
// IMPORTANT function.
func exeTmpl(w http.ResponseWriter, r *http.Request, view *viewData, tmpl string) {
	if view == nil {
		view = &viewData{Credentials: &credentials{User: &user{}}}

	}
	if r.Context().Value(ctxkey) != nil {
		view.Credentials = r.Context().Value(ctxkey).(*credentials)
	}
	// if view.Profile == nil {
	// 	view.Credentials = r.Context().Value(ctxkey).(*credentials)
	// }

	view.AppName = AppName
	view.Stream = stream
	err := templates.ExecuteTemplate(w, tmpl, view)
	if err != nil {
		log.Println(err)
	}
}

// ajaxResponse()is used to respond to ajax requests with arbitrary data in the
// format of map[string]string
func ajaxResponse(w http.ResponseWriter, res map[string]string) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Println(err)
	}
}

// status() is kind of a wrapper around ajaxResponse, allowing us to send a
// generic response, with only a "status" attribute set to s, which is usually
// either "success", or some error. USE:
// log.Println(status(w, "Error Summary", err))
//
//	or
//
// log.Println(status(w, "success", nil))
func status(w http.ResponseWriter, s string, err error) error {
	ajaxResponse(w, map[string]string{"status": s})
	return err
}

// genID(length int) generates an item ID of length length by picking s random
// character from a string containing the upper and lower alphabet and numbers
// 0-9. We omit the characters E,F,G, and H to keep the line containing the
// symbols variable from exceeding 80 characters, and to reduce the codebase by
// abstaining from adding a line break.
func genID(length int) (ID string) {
	symbols := "abcdefghijklmnopqrstuvwxyz1234567890ABCDIJKLMNOPQRSTUVWXYZ"
	for i := 0; i <= length; i++ {
		s := rand.Intn(len(symbols))
		ID += symbols[s : s+1]
	}
	return
}

// makeZmem() returns a redis Z member for use in a ZSET. Score is set to zero.
func makeZmem(st string) redis.Z {
	return redis.Z{Member: st, Score: 0}
}

// marshalPostData() is used to marshal a post in the form of a JSON string
// sent by the client into a *post{} struct.
func marshalPostData(r *http.Request) (*post, error) {
	t := &post{}
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	err := decoder.Decode(t)
	if err != nil {
		return t, err
	}
	return t, nil
}

// setupLogging() sets output flags and the file for logging, allowing for the
// logging of line numbers of errors (useful!). Also sets the logfile file.
func setupLogging() {
	if len(logFilePath) > 1 {
		f, err := os.OpenFile(
			logFilePath, // see: bolt.conf.json
			os.O_RDWR|os.O_CREATE|os.O_APPEND,
			0666)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		log.SetOutput(f)
		defer f.Close()
	}
}

// readConf() is used to read the bolt.conf.json configuration file, which is
// used to configure the port, log file, and app name, but is otherwise
// unnecessary for the tagmachine app to function.
func readConf() *config {
	b, err := os.ReadFile("./bolt.conf.json")
	if err != nil {
		log.Println(err)
	}
	c := config{}
	err = json.Unmarshal(b, &c)
	if err != nil {
		log.Println(err)
	}
	return &c
}
