// main_handlers.go houses route handlers for tagmachine which aren't long
// enough to warrant their own file. Handlers located outside of
// main_handlers.go will be suffixed with _handler.go. We try to keep the
// handler functions in some logical order.
package main

import (
	"fmt"
	"log"
	"net/http"
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

// viewItem() is the route handler used for viewing a link to an individual
// post. It serves "main.html", passing the single post as the "stream" value
// in viewData{}, (allowing us to reuse "main.html", instead of creating
// another page view).
func viewItem(w http.ResponseWriter, r *http.Request) {
	// get the ID from after the "view/", the route looks like this:
	// https://tagmachine.xyz/view/LGnIKd2DXECZPsBQ
	id := strings.Split(r.RequestURI, "/")[2]

	// Execute the template with the single post added as the
	// viewData.Stream{} property.
	exeTmpl(w, r, &viewData{
		AppName: appConf.App.Name,
		Stream:  getPostsByID([]string{id}),
	}, "main.html")
}

// profileHandler() is the route handler used for viewing a users profile. It
// parses the ID from the request URI, serving the profile associated with that
// user ID, using the page view "profile.html",
func profileHandler(w http.ResponseWriter, r *http.Request) {
	// get the ID from after "user/", the route looks like this:
	// https://tagmachine.xyz/user/LGnIKd2DXECZPsBQ
	id := strings.Split(r.RequestURI, "/")[2]

	// Create a dummy credentials{} with the ID for credentials.User set
	// with the ID obtained above.
	var _c *credentials = &credentials{
		User: &user{ID: id},
	}

	// Get the profile data for the user by passing the dummy credentials
	// to scanProfile().
	err := scanProfile(_c) // see: scanProfile()
	if err != nil {
		log.Println(status(w, "Couldn't find user", err))
		return
	}

	// If user.ProfilePic is unset we give it a default value.
	if _c.User.ProfilePic == "" {
		_c.User.ProfilePic = "public/media/ndt.jpg"
	}

	// If user.ProfileBG is unset we give it a default value.
	if _c.User.ProfileBG == "" {
		_c.User.ProfileBG = "public/media/hubble.jpg"
	}

	// Get the users liked posts to show visitors to their profile.
	likes, err := getLikes(_c) // see: getLikes()
	if err != nil {
		log.Println(status(w, "Database error", err))
		return
	}

	// Execute the "profile.html" page view template with the dummy users
	// profile information set as the viewData{}.Profile property,
	// providing the authenticated users credentials and the stream of
	// liked posts associated with the profile being viewed as well.
	exeTmpl(w, r, &viewData{
		Profile:     _c.User,
		Credentials: r.Context().Value(ctxkey).(*credentials),
		Stream:      likes,
	}, "profile.html")
}

// reply() is the route handler for post replies, and is an ajax
// response/request, thus we don't redirect the client to a new page, but may
// return data to update the view.
func reply(w http.ResponseWriter, r *http.Request) {
	// marshal the post data, checking for integrity and validity:
	p, err := marshalPostData(r)
	if err != nil {
		log.Println(status(w, "Invalid Data?", err))
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
		log.Println(status(w, "Database Error", err))
		return
	}

	// success
	ajaxResponse(w, map[string]string{"status": "success", "ID": p.ID})
}

// likeHandler() is the route handler for /like/ID, by appending the liked
// posts ID to the the users []user.Likes slice. and is triggered when a user
// likes or unlikes a post.
func likeHandler(w http.ResponseWriter, r *http.Request) {
	// parse the URI for the ID of the post being liked/unliked.
	id := strings.Split(r.RequestURI, "/")[2]

	// Get the user object from the context.
	c := r.Context().Value(ctxkey).(*credentials)

	// Update the database.
	n, err := setLike(c, id) // see: setLike()
	if err != nil {
		log.Println(status(w, "Database Error", err))
		return
	}

	// Append the liked posts ID to the users []user.Likes slice.
	c.User.Likes = append(c.User.Likes, id)

	// Save the updates to the user to the database.
	err = setProfile(c) // see: setProfile()
	if err != nil {
		log.Println(status(w, "Database Error", err))
		return
	}

	// success
	ajaxResponse(w, map[string]string{
		"success": "true",
		"score":   fmt.Sprint(1 - n),
	})
}

// addFriendHandler() is used to add a friend to a users friends list, by
// appending the friends ID to the the users []user.Friends slice.
func addFriendHandler(w http.ResponseWriter, r *http.Request) {
	// get the ID by parsing the request URI.
	id := strings.Split(r.RequestURI, "/")[2]

	// Get user object from the context.
	c := r.Context().Value(ctxkey).(*credentials)

	// TODO:
	n, err := setFriend(c, id) // see: setFriend()
	if err != nil {
		log.Println(status(w, "Database Error", err))
		return
	}

	// Append the new freiends ID to the users []user.Friends slice, so
	// that we can look them up later.
	c.User.Friends = append(c.User.Friends, id)

	// Save the changes to the users profile.
	err = setProfile(c) // see: setProfile()
	if err != nil {
		log.Println(err)
	}

	// success. We send back the score.
	// TODO: Re-implement.
	ajaxResponse(w, map[string]string{
		"success": "true",
		"score":   fmt.Sprint(1 - n),
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
func shareHandler(w http.ResponseWriter, r *http.Request) {
	exeTmpl(w, r, nil, "main.html")
}
func unFriendHandler(w http.ResponseWriter, r *http.Request) {
	exeTmpl(w, r, nil, "main.html")
}

// func editHandler(w http.ResponseWriter, r *http.Request) {
// 	post, err := parseForm(r)
// 	if err != nil {
// 		log.Println(status(w, "Invalid Form", err))
// 		return
// 	}
// 	var c *credentials = r.Context().Value(ctxkey).(*credentials)
// 	log.Println(post, c, post.TempFileName)
// 	switch post.Type {
// 	case "ProfilePic":
// 		c.User.ProfilePic = post.TempFileName
// 	case "ProfileBG":
// 		c.User.ProfileBG = post.TempFileName
// 	}
// 	log.Println(c.User, c)
// 	// c.User.Posts = nil
// 	err = setProfile(c)
// 	if err != nil {
// 		log.Println(status(w, "Database Error", err))
// 		return
// 	}
// 	ajaxResponse(w, map[string]string{
// 		"status":  "success",
// 		"payload": post.TempFileName,
// 	})
// }

// likesHandler() is the route handler for /user/ID/likes
//
//	func likesHandler(w http.ResponseWriter, r *http.Request) {
//		id := strings.Split(r.RequestURI, "/")[2]
//		c := r.Context().Value(ctxkey).(*credentials)
//		n, err := setLike(c, id)
//		if err != nil {
//			log.Println(err)
//			return
//		}
//		c.User.Likes = append(c.User.Likes, id)
//		err = setProfile(c)
//		if err != nil {
//			log.Println(err)
//		}
//		ajaxResponse(w, map[string]string{
//                      "success": "true",
//                      "score": fmt.Sprint(1 - n),
//              })
//	}
