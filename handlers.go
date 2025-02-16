package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

func root(w http.ResponseWriter, r *http.Request) {
	exeTmpl(w, r, nil, "main.html")
}
func viewItem(w http.ResponseWriter, r *http.Request) {
	id := strings.Split(r.RequestURI, "/")[2]
	exeTmpl(w, r, &viewData{
		AppName: appConf.App.Name,
		Stream:  getPostsByID([]string{id}),
	}, "main.html")
}

func likeHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.Split(r.RequestURI, "/")[2]
	c := r.Context().Value(ctxkey).(*credentials)
	n, err := setLike(c, id)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(c)
	c.User.Likes = append(c.User.Likes, id)
	err = setProfile(c)
	if err != nil {
		log.Println(err)
	}
	log.Println(id, 0, n)
	ajaxResponse(w, map[string]string{"success": "true", "score": fmt.Sprint(1 - n)})
}
func likesHandler(w http.ResponseWriter, r *http.Request) {
	exeTmpl(w, r, nil, "main.html")
}
func shareHandler(w http.ResponseWriter, r *http.Request) {
	exeTmpl(w, r, nil, "main.html")
}
func addFriendHandler(w http.ResponseWriter, r *http.Request) {
	exeTmpl(w, r, nil, "main.html")
}
func unFriendHandler(w http.ResponseWriter, r *http.Request) {
	exeTmpl(w, r, nil, "main.html")
}
func profileHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.Split(r.RequestURI, "/")[2]
	log.Println(id, "------")
	var _c *credentials = &credentials{
		User: &user{ID: id},
	}

	err := scanProfile(_c)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(_c.User)
	if _c.User.ProfilePic == "" {
		_c.User.ProfilePic = "public/media/ndt.jpg"
	}
	if _c.User.ProfileBG == "" {
		_c.User.ProfileBG = "public/media/hubble.jpg"
	}

	exeTmpl(w, r, &viewData{Profile: _c.User, Credentials: r.Context().Value(ctxkey).(*credentials)}, "profile.html")
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	post, err := parseForm(r)
	if err != nil {
		log.Println(status(w, "Invalid Form", err))
		return
	}
	var c *credentials = r.Context().Value(ctxkey).(*credentials)
	switch post.Type {
	case "ProfilePic":
		c.User.ProfilePic = post.TempFileName
	case "ProfileBG":
		c.User.ProfileBG = post.TempFileName
	}
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
