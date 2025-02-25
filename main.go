// Provided Under BSD (2 Clause)
//
// Copyright 2025 Johnathan A. Hartsfield
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice,
// this list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice,
// this list of conditions and the following disclaimer in the documentation
// and/or other materials provided with the distribution.
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
// upload_handler.go houses the functions used for user post submissions and
// other multipart/form-data requests. This file may be split up in the future.
// main.go houses the main() and init() functions, and struct definitions,
// along with any necessary implementations (encoding.BinaryMarshaler is
// needed to marshal data for redis).
//
// ////////////////////////////////////////////////////////////////////////////
// ////////////////////////////////////////////////////////////////////////////
// /////////////////                                         //////////////////
// /////////////////                                         //////////////////
// /////////////////                                         //////////////////
// /////////////////       <- #[Tag]!Machine( $%^@~ )        //////////////////
// /////////////////                                         //////////////////
// /////////////////                                         //////////////////
// /////////////////                                         //////////////////
// ////////////////////////////////////////////////////////////////////////////
// ////////////////////////////////////////////////////////////////////////////
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/redis/go-redis/v9"
)

// ckey/ctxkey is used as the key for the HTML context and is how we retrieve
// token information and pass it around to handlers.
type ckey int

// used as a type of "nonce" within our http context.
const ctxkey ckey = iota

var (
	// Used for decoding the user token and should be provided as an
	// environment variable at run time for security:
	hmacSampleSecret []byte = []byte(os.Getenv("hmacss"))
	// read the bolt.conf.json file and obtain some basic configuration
	// variables such as appName, port, and logFilePath.
	appConf *config = readConf()
	// AppName is used in some templates.
	AppName     string = appConf.App.Name
	logFilePath string = appConf.App.Env["logFilePath"]
	servicePort string = ":" + appConf.App.Port

	// initialize templates.
	templates *template.Template = template.New("")

	// initialize post stream.
	stream []*post = []*post{}

	// Create the context for redis, and connect tot he redis database.
	rdx context.Context = context.Background()
	rdb *redis.Client   = redis.NewClient(&redis.Options{
		Addr:     ":6379",
		Password: "",
		DB:       2,
	})
)

// config{} is used by readConf() to read the bolt.conf.json file.
type config struct {
	App struct {
		Name       string            `json:"name" redis:"name"`
		DomainName string            `json:"domain_name" redis:"domain_name"`
		Version    string            `json:"version" redis:"version"`
		Env        map[string]string `json:"env" redis:"env"`
		Port       string            `json:"port" redis:"port"`
		AlertsOn   bool              `json:"alertsOn" redis:"alerts_on"`
		TLSEnabled bool              `json:"tls_enabled" redis:"tls_enabled"`
		Repo       string            `json:"repo" redis:"repo"`
	} `json:"app" redis:"app"`
	GCloud struct {
		Command   string `json:"command" redis:"command"`
		Zone      string `json:"zone" redis:"zone"`
		Project   string `json:"project" redis:"project"`
		User      string `json:"user" redis:"user"`
		LiveDir   string `json:"livedir" redis:"live_dir"`
		ProxyConf string `json:"proxyConf" redis:"proxy_conf"`
	} `json:"gcloud" redis:"g_cloud"`
}

// viewData{} represents the root model used to dynamically update the page
// views, and is passed to the client with each page request, but not typically
// in ajax responses.
type viewData struct {
	// AppName is the AppName as found in bolt.conf.json.
	AppName string `json:"app_name" redis:"app_name"`
	// Stream is a stream of posts, which could be one or many, in multiple
	// sort orders.
	Stream []*post `json:"stream" redis:"stream"`
	// Nonce is a number used once, which helps prevent double posting and
	// (in theory) mitigates certain attack vectors.
	Nonce int `json:"nonce" redis:"nonce"`
	// Credentials are used for logging a user in and contain a logged in
	// users credentials.
	Credentials *credentials `json:"credentials" redis:"credentials"`
	// Profile is used when viewing another users profile (or when a user
	// views their own profile.)
	Profile *user `json:"user" redis:"user"`
}

// credentials are user credentials and are used in the HTML templates and also
// by handlers that do authorized requests.
type credentials struct {
	// Name is used for login, and will be set to the users login email.
	Name string `json:"username" redis:"username"`
	// Password is used for login, transferred over SSL, and never stored
	// in plain text.
	Password string `json:"password" redis:"password"`
	// Used in templates to determine whether a user is logged on.
	IsLoggedIn bool `json:"isLoggedIn" redis:"isLoggedIn"`
	// The logged in user, if any. Sometimes a "dummy" user with a user ID
	// is placed here to look up user data.
	User *user `json:"user" redis:"user"`
	// Implements
	jwt.StandardClaims
}

// *credentials.UnmarshalBinary() is used to implement
// encoding.BinaryMarshaler, as required for compatibility with redis.
func (u *credentials) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, u)
}

// *credentials.MarshalBinary() is used to implement encoding.BinaryMarshaler,
// as required for compatibility with redis.
func (p *credentials) MarshalBinary() ([]byte, error) {
	return json.Marshal(p)
}

// post{} represents a user post or reply to another users post.
type post struct {
	Type         string        `json:"Type" redis:"Type"`
	ID           string        `json:"id" redis:"id"`
	Parent       string        `json:"parent" redis:"parent"`
	TS           time.Time     `json:"ts" redis:"ts"`
	TimeString   string        `json:"time_string" redis:"time_string"`
	Author       string        `json:"author" redis:"author"`
	Text         template.HTML `json:"uptext" redis:"uptext"`
	Media        string        `json:"Media" redis:"Media"`
	MediaType    template.HTML `json:"media_type" redis:"media_type"`
	TempFileName string        `json:"temp_file_name" redis:"temp_file_name"`
	Score        int           `json:"score" redis:"score"`
	Categories   []string      `json:"categories" redis:"categories"`
	Tags         []*tag        `json:"tags" redis:"tags"`
	CommentIDs   []string      `json:"commentIDs" redis:"commentIDs"`
	Comments     []*post       `json:"comments" redis:"comments"`
}

// *post.UnmarshalBinary() is used to implement encoding.BinaryMarshaler, as
// required for compatibility with redis.
func (u *post) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, u)
}

// *post.MarshalBinary() is used to implement encoding.BinaryMarshaler, as
// required for compatibility with redis.
func (p *post) MarshalBinary() ([]byte, error) {
	return json.Marshal(p)
}

// user{} represents a user.
type user struct {
	// Token is where we store the users token locally, and is generally
	// retrieved from redis by looking up the users hash map via their ID
	//
	Token       string    `json:"token" redis:"token"`
	ID          string    `json:"id" redis:"id"`
	Email       string    `json:"email" redis:"email"`
	Score       int       `json:"score" redis:"score"`
	Joined      time.Time `json:"joined" redis:"joined"`
	LastSeen    time.Time `json:"last_seen" redis:"last_seen"`
	About       string    `json:"about" redis:"about"`
	Work        string    `json:"work" redis:"work"`
	Location    string    `json:"location" redis:"location"`
	ProfileBG   string    `json:"profile_bg" redis:"profile_bg"`
	ProfilePic  string    `json:"profile_pic" redis:"profile_pic"`
	ProfilePics []string  `json:"profile_pics" redis:"profile_pics"`
	Posts       []string  `json:"posts" redis:"posts"`
	Likes       []string  `json:"likes" redis:"likes"`
	Shares      []string  `json:"shares" redis:"shares"`
	Friends     []string  `json:"friends" redis:"friends"`

	// TODO /* Not implemented */
	Events   []string `json:"events" redis:"events"`
	Level    int      `json:"level" redis:"level"`
	Throttle int      `json:"throttle" redis:"throttle"`
	Insights string   `json:"insights" redis:"insights"`
	Random   string   `json:"random" redis:"random"`
	Other    string   `json:"other" redis:"other"`
}

// *user.UnmarshalBinary() is used to implement encoding.BinaryMarshaler, as
// required for compatibility with redis.
func (u *user) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, u)
}

// *user.MarshalBinary() is used to implement encoding.BinaryMarshaler, as
// required for compatibility with redis.
func (u *user) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

// tag is a tag which is used to categorize and possibly analyze a post.
type tag struct {
	// ID is the tag name basically.
	ID string `json:"id" redis:"id"`
	// Count is the number of times the tags been used.
	Count int `json:"count" redis:"count"`
	// Score is the tags score, which is determined by a YTBD algorithm.
	Score int `json:"score" redis:"score"`
	// Born is the date of the first occurrence of this tag.
	Born time.Time `json:"born" redis:"born"`
}

// *tag.UnmarshalBinary() is used to implement encoding.BinaryMarshaler, as
// required for compatibility with redis.
func (u *tag) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, u)
}

// *tag.MarshalBinary() is used to implement encoding.BinaryMarshaler, as
// required for compatibility with redis.
func (p *tag) MarshalBinary() ([]byte, error) {
	return json.Marshal(p)
}

// init() sets some logging flags to allow display of line numbers, parses the
// templates contained in internal/, and makes the ./public/temp/ dir if it
// hasn't already been created.
func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	template.Must(templates.ParseGlob("internal/*/*/*"))
	err := os.Mkdir("./public/temp", 0777)
	if err != nil {
		log.Println(err)
	}
}

// main() sets up logging by initializing it, begins an infinite loop in a go
// function, used for caching the database every two seconds (will be
// reconfigured for production).
func main() {
	setupLogging()
	go func() {
		for {
			cache()
			time.Sleep(2 * time.Second)
		}
	}()

	// start the server.
	ctx, srv := bolt()

	// print the server address to the terminal.
	log.Println("@ http://localhost" + srv.Addr)
	fmt.Println("\n\n@ http://localhost" + srv.Addr +
		"  -->  " + appConf.App.DomainName)

	// hold the program open until done.
	<-ctx.Done()
}
