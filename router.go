package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"
)

func bolt() (ctx context.Context, srv *http.Server) {
	var mux *http.ServeMux = http.NewServeMux()
	registerRoutes(mux)
	mux.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
	srv = serverFromConf(mux)
	ctx, cancelCtx := context.WithCancel(context.Background())
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Panicln(err)
		}
		cancelCtx()
	}()
	return
}

// serverFromConf returns a *http.Server with a pre-defined configuration
func serverFromConf(mux *http.ServeMux) *http.Server {
	return &http.Server{
		Addr:              servicePort,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       5 * time.Second,
	}
}

var lastSave string

func wasmodified(w http.ResponseWriter, r *http.Request) {
	b, err := os.ReadFile(".lastsavetime_bolt")
	if err != nil {
		log.Println(err)
	}
	if string(b) != lastSave {
		lastSave = string(b)
		ajaxResponse(w, map[string]string{"modified": "true"})
		return
	}
	ajaxResponse(w, map[string]string{"modified": "false"})
}

// registerRoutes registers the routes with the provided *http.ServeMux
func registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", checkAuth(root))
	mux.HandleFunc("/reply", checkAuth(reply))
	mux.HandleFunc("/what", what)
	mux.HandleFunc("/signin", signin)
	mux.HandleFunc("/signup", signup)
	mux.HandleFunc("/signout", signout)
	mux.HandleFunc("/uploadItem", checkAuth(uploadHandler))
	mux.HandleFunc("/view/", checkAuth(viewItem))
	mux.HandleFunc("/like/", checkAuth(likeHandler))
	mux.HandleFunc("/share/", checkAuth(shareHandler))
	mux.HandleFunc("/addFriend/", checkAuth(addFriendHandler))
	mux.HandleFunc("/unfriend/", checkAuth(unFriendHandler))
	mux.HandleFunc("/edit", checkAuth(editHandler))
	// mux.HandleFunc("/save", checkAuth(editProfileHandler))
	mux.HandleFunc("/tag/", checkAuth(tagHandler))
	mux.HandleFunc("/friends/", friendHandler)
	mux.HandleFunc("/search/", searchHandler)
	mux.HandleFunc("/user/", checkAuth(profileHandler))
	mux.HandleFunc("/likes/", likesHandler)
	mux.HandleFunc("/wasmodified", wasmodified)
}
