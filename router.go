// / Provided Under BSD (2 Clause)
//
// Copyright 2025 Johnathan A. Hartsfield
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
//  1. Redistributions of source code must retain the above copyright notice,
//     this list of conditions and the following disclaimer.
//
//  2. Redistributions in binary form must reproduce the above copyright notice,
//     this list of conditions and the following disclaimer in the documentation
//     and/or other materials provided with the distribution.
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
