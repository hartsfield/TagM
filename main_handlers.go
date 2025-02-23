package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// root() is the route handler for the "home page" of tagmachine, which is what
// a visitor will see when they visit tagmachine.xyz. It serves "main.html",
// passing nil, as the viewData, which will cause viewData to be set to default
// values determined by exeTmpl.
func root(w http.ResponseWriter, r *http.Request) {
	exeTmpl(w, r, nil, "main.html")
}

// what() is the route handler for tagmachine.xyz/what, which is basically
// the about page. It serves "what.html", passing nil, as the viewData, which
// will cause viewData to be set to default values determined by exeTmpl.
func what(w http.ResponseWriter, r *http.Request) {
	exeTmpl(w, r, nil, "what.html")
}

// reply() is the route handler for post replies, and is an ajax
// response/request, thus we don't redirect the client to a new page, but may
// return data to update the view.
func reply(w http.ResponseWriter, r *http.Request) {
	// marshal the post data, checking for integrity and validity:
	p, err := marshalPostData(r)
	if err != nil {
		log.Println(err)
		ajaxResponse(w, map[string]string{"status": err.Error()})
		return
	}

	// get the users credentials from the context
	var c_ *credentials = r.Context().Value(ctxkey).(*credentials)

	// initialize some default variables for the new reply (which is a
	// post{}).
	p.ID = genID(15)
	p.TS = time.Now()
	p.TimeString = time.Now().Format(time.RFC822)
	p.Author = c_.User.ID

	// Add the posts ID to a sorted set and store the post data as an
	// object in redis:
	_, err = zaddUsersPosts(c_, p)
	if err != nil {
		log.Println(err)
		ajaxResponse(w, map[string]string{"status": err.Error()})
		return
	}

	// success
	ajaxResponse(w, map[string]string{"status": "success", "ID": p.ID})
}

// viewItem() is the route handler used for viewing a link to an individual
// post. It serves "main.html", passing the single post as the "stream" value
// in viewData{}, (allowing us to reuse "main.html", instead of creating
// another page view).
func viewItem(w http.ResponseWriter, r *http.Request) {
	// get the ID from after the "view/", the route looks like this:
	// https://tagmachine.xyz/view/LGnIKd2DXECZPsBQ
	id := strings.Split(r.RequestURI, "/")[2]

	exeTmpl(w, r, &viewData{
		AppName: appConf.App.Name,
		Stream:  getPostsByID([]string{id}),
	}, "main.html")
}

// likeHandler is the routehandler for /like/ID, and is used when a user likes
// or unlikes a post.
func likeHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.Split(r.RequestURI, "/")[2]
	c := r.Context().Value(ctxkey).(*credentials)
	n, err := setLike(c, id)
	if err != nil {
		log.Println(err)
		return
	}
	c.User.Likes = append(c.User.Likes, id)
	err = setProfile(c)
	if err != nil {
		log.Println(err)
	}
	ajaxResponse(w, map[string]string{"success": "true", "score": fmt.Sprint(1 - n)})
}
func likesHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.Split(r.RequestURI, "/")[2]
	c := r.Context().Value(ctxkey).(*credentials)
	n, err := setLike(c, id)
	if err != nil {
		log.Println(err)
		return
	}
	c.User.Likes = append(c.User.Likes, id)
	err = setProfile(c)
	if err != nil {
		log.Println(err)
	}
	ajaxResponse(w, map[string]string{"success": "true", "score": fmt.Sprint(1 - n)})
}
func shareHandler(w http.ResponseWriter, r *http.Request) {
	exeTmpl(w, r, nil, "main.html")
}
func addFriendHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.Split(r.RequestURI, "/")[2]
	c := r.Context().Value(ctxkey).(*credentials)
	n, err := setFriend(c, id)
	if err != nil {
		log.Println(err)
		return
	}
	c.User.Friends = append(c.User.Friends, id)
	err = setProfile(c)
	if err != nil {
		log.Println(err)
	}
	ajaxResponse(w, map[string]string{"success": "true", "score": fmt.Sprint(1 - n)})
}
func unFriendHandler(w http.ResponseWriter, r *http.Request) {
	exeTmpl(w, r, nil, "main.html")
}
func profileHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.Split(r.RequestURI, "/")[2]
	var _c *credentials = &credentials{
		User: &user{ID: id},
	}

	err := scanProfile(_c)
	if err != nil {
		log.Println(err)
	}
	if _c.User.ProfilePic == "" {
		_c.User.ProfilePic = "public/media/ndt.jpg"
	}
	if _c.User.ProfileBG == "" {
		_c.User.ProfileBG = "public/media/hubble.jpg"
	}
	likes, err := getLikes(_c)
	if err != nil {
		log.Println(err)
	}
	exeTmpl(w, r, &viewData{
		Profile:     _c.User,
		Credentials: r.Context().Value(ctxkey).(*credentials),
		Stream:      likes,
	}, "profile.html")
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	post, err := parseForm(r)
	if err != nil {
		log.Println(status(w, "Invalid Form", err))
		return
	}
	var c *credentials = r.Context().Value(ctxkey).(*credentials)
	log.Println(post, c, post.TempFileName)
	switch post.Type {
	case "ProfilePic":
		c.User.ProfilePic = post.TempFileName
	case "ProfileBG":
		c.User.ProfileBG = post.TempFileName
	}
	log.Println(c.User, c)
	// c.User.Posts = nil
	err = setProfile(c)
	if err != nil {
		log.Println(status(w, "Database Error", err))
		return
	}
	ajaxResponse(w, map[string]string{
		"status":  "success",
		"payload": post.TempFileName,
	})
}
func friendHandler(w http.ResponseWriter, r *http.Request) {
	exeTmpl(w, r, nil, "main.html")
}
func searchHandler(w http.ResponseWriter, r *http.Request) {
	exeTmpl(w, r, nil, "main.html")
}
func tagHandler(w http.ResponseWriter, r *http.Request) {
	exeTmpl(w, r, nil, "main.html")
}

// /////////////////////////////////////////////////////////////////////////////
// auto reload (development)
// /////////////////////////////////////////////////////////////////////////////
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

///////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////
