package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"

	"github.com/redis/go-redis/v9"
)

// exeTmpl is used to build and execute an html template.
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

// makeZmem returns a redis Z member for use in a ZSET. Score is set to zero
func makeZmem(st string) redis.Z {
	return redis.Z{
		Member: st,
		Score:  0,
	}
}

// ajaxResponse is used to respond to ajax requests with arbitrary data in the
// format of map[string]string
func ajaxResponse(w http.ResponseWriter, res map[string]string) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Println(err)
	}
}

// genID generates a item ID
func genID(length int) (ID string) {
	symbols := "abcdefghijklmnopqrstuvwxyz1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	for i := 0; i <= length; i++ {
		s := rand.Intn(len(symbols))
		ID += symbols[s : s+1]
	}
	return
}

// setupLogging sets output flags and the file for logging
func setupLogging() {
	if len(logFilePath) > 1 {

		f, err := os.OpenFile(logFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}

		log.SetOutput(f)
		defer f.Close()
	}
}

// readConf is used to read the bolt.conf.json configuration file.
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

// marshalPostData is used to marshal a JSON string into a *post struct
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
