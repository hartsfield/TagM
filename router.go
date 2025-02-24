// router.go is where the tagmachines routes are defined, and also where the
// server is configured and initialized.
package main

import (
	"context"
	"log"
	"net/http"
	"time"
)

// bolt() starts the http(s) server.
func bolt() (ctx context.Context, srv *http.Server) {
	var mux *http.ServeMux = http.NewServeMux()
	registerRoutes(mux)

	// Tell the server /public is accessible to the world wide web.
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

// serverFromConf() returns a *http.Server with a pre-defined configuration
// using the multiplexer at the bottom of this file.
func serverFromConf(mux *http.ServeMux) *http.Server {
	return &http.Server{
		// servicePort is in main.go, and configured in bolt.json.conf.
		Addr:              servicePort,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       5 * time.Second,
	}
}

// registerRoutes registers the routes with the provided *http.ServeMux. Add
// new routes here as necessary, wrapping the route handler with the
// checkAuth(handler) middle ware function as needed. Keep the multiplexer at
// the bottom of this file to allow for the programmatic insertion of routes
// via external tools.
func registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", checkAuth(root))
	mux.HandleFunc("/reply", checkAuth(reply))
	mux.HandleFunc("/what", what)
	mux.HandleFunc("/signin", signin)
	mux.HandleFunc("/signup", signup)
	mux.HandleFunc("/uploadItem", checkAuth(uploadHandler))
	mux.HandleFunc("/view/", checkAuth(viewItem))
	mux.HandleFunc("/like/", checkAuth(likeHandler))
	mux.HandleFunc("/share/", checkAuth(shareHandler))
	mux.HandleFunc("/addFriend/", checkAuth(addFriendHandler))
	mux.HandleFunc("/unfriend/", checkAuth(unFriendHandler))
	mux.HandleFunc("/tag/", checkAuth(tagHandler))
	mux.HandleFunc("/friends/", friendHandler)
	mux.HandleFunc("/search/", searchHandler)
	mux.HandleFunc("/user/", checkAuth(profileHandler))
	// mux.HandleFunc("/edit", checkAuth(editHandler))
	// mux.HandleFunc("/likes/", likesHandler)
}
