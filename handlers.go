package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

func root(w http.ResponseWriter, r *http.Request) {
	exeTmpl(w, r, nil, "main.html")
}
func what(w http.ResponseWriter, r *http.Request) {
	exeTmpl(w, r, nil, "what.html")
}
func reply(w http.ResponseWriter, r *http.Request) {
	p, err := marshalPostData(r)
	if err != nil {
		log.Println(err)
	}
	var c_ *credentials = r.Context().Value(ctxkey).(*credentials)
	p.ID = genID(15)
	p.TS = time.Now()
	p.TimeString = time.Now().Format(time.RFC822)
	p.Author = c_.User.ID
	_, err = zaddUsersPosts(c_, p)
	if err != nil {
		log.Println(err)
	}
	ajaxResponse(w, map[string]string{})
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
