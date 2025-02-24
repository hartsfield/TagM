// main.go houses the main() and init() functions, and struct definitions,
// along with any necessary implementations (encoding.BinaryMarshaler is
// needed to marshal data for redis).
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
// token information and pass it around to handlers
type ckey int

const ctxkey ckey = iota

var (
	// Used for decoding the user token and should be provided as an
	// environment variable at run time for security:
	hmacSampleSecret []byte = []byte(os.Getenv("hmacss"))
	// read the bolt.conf.json file and obtain some basic configuration
	// variables such as appName, port, and logFilePath.
	appConf     *config = readConf()
	AppName     string  = appConf.App.Name
	logFilePath string  = appConf.App.Env["logFilePath"]
	servicePort string  = ":" + appConf.App.Port

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

// viewData{} represents the root model used to dynamically update the page views.
type viewData struct {
	AppName     string       `json:"app_name" redis:"app_name"`
	Stream      []*post      `json:"stream" redis:"stream"`
	Nonce       int          `json:"nonce" redis:"nonce"`
	Credentials *credentials `json:"credentials" redis:"credentials"`
	Profile     *user        `json:"user" redis:"user"`
}

// post{} represents a user post or reply to another users post.
type post struct {
	Parent       string        `json:"parent" redis:"parent"`
	Categories   []string      `json:"categories" redis:"categories"`
	Media        string        `json:"Media" redis:"Media"`
	Type         string        `json:"Type" redis:"Type"`
	Author       string        `json:"author" redis:"author"`
	Text         template.HTML `json:"uptext" redis:"uptext"`
	ID           string        `json:"id" redis:"id"`
	TS           time.Time     `json:"ts" redis:"ts"`
	TimeString   string        `json:"time_string" redis:"time_string"`
	MediaType    template.HTML `json:"media_type" redis:"media_type"`
	TempFileName string        `json:"temp_file_name" redis:"temp_file_name"`
	Tags         []*tag        `json:"tags" redis:"tags"`
	Score        int           `json:"score" redis:"score"`
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

type tag struct {
	Tag   string    `json:"tag" redis:"tag"`
	Count int       `json:"count" redis:"count"`
	Born  time.Time `json:"born" redis:"born"`
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

// credentials are user credentials and are used in the HTML templates and also
// by handlers that do authorized requests
type credentials struct {
	Name       string `json:"username" redis:"username"`
	Password   string `json:"password" redis:"password"`
	IsLoggedIn bool   `json:"isLoggedIn" redis:"isLoggedIn"`
	User       *user  `json:"user" redis:"user"`
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

// user{} represents a user.
type user struct {
	Token       string    `json:"token" redis:"token"`
	ID          string    `json:"id" redis:"id"`
	Email       string    `json:"email" redis:"email"`
	ProfilePic  string    `json:"profile_pic" redis:"profile_pic"`
	ProfilePics []string  `json:"profile_pics" redis:"profile_pics"`
	ProfileBG   string    `json:"profile_bg" redis:"profile_bg"`
	Location    string    `json:"location" redis:"location"`
	About       string    `json:"about" redis:"about"`
	Work        string    `json:"work" redis:"work"`
	Score       int       `json:"score" redis:"score"`
	Posts       []string  `json:"posts" redis:"posts"`
	Likes       []string  `json:"likes" redis:"likes"`
	Shares      []string  `json:"shares" redis:"shares"`
	Friends     []string  `json:"friends" redis:"friends"`
	Events      []string  `json:"events" redis:"events"`
	Joined      time.Time `json:"joined" redis:"joined"`
	LastSeen    time.Time `json:"last_seen" redis:"last_seen"`
	Insights    string    `json:"insights" redis:"insights"`
	Level       int       `json:"level" redis:"level"`
	Throttle    int       `json:"throttle" redis:"throttle"`
	Random      string    `json:"random" redis:"random"`
	Other       string    `json:"other" redis:"other"`
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
	fmt.Println("\n\n@ http://localhost" + srv.Addr + "  -->  " + appConf.App.DomainName)

	// hold the program open until done.
	<-ctx.Done()
}
