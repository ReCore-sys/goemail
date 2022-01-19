package main

import (
	"database/sql"

	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	sqlstuff "github.com/ReCore-sys/gomail/libraries/sql"
	utils "github.com/ReCore-sys/gomail/libraries/utils"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3" // We null this cos idfk
)

const (
	domain     = "mommymilkies.org"
	prettyname = "Mommy Milkies"
)

var filelocation, _ = os.Getwd()

// Getdircontent gets a slice of all the files in a dir
func Getdircontent(dir string) (content []string) {
	content = make([]string, 0)
	files, _ := ioutil.ReadDir(dir)

	for _, f := range files {
		content = append(content, f.Name())
	}
	return content
}

type news struct {
	Title   string `json:"title"`
	Date    string `json:"date"`
	Content string `json:"content"`
	Update  string `json:"update"`
}

// Letsa go!
func main() {

	// Create the database
	if utils.CheckFile("data.db") == false {
		os.Create(filelocation + "/data.db")
	}
	conn, err := sql.Open("sqlite3", filelocation+"/data.db") // Open the database
	if err != nil {
		log.Fatal(err) // If there is an error, log it
	}
	db := sqlstuff.SQL{Sqlcon: conn, Tablename: "usr"}
	// Set up gin
	r := gin.Default()
	r.SetTrustedProxies(nil) // If I don't do this, gin will complain about untrusted proxies. I don't know why.
	// Set up paths to the templates and the static files
	r.LoadHTMLGlob("templates/*.html")
	r.Static("/css", "static/css")
	r.Static("/js", "static/js")
	r.Static("/img", "static/img")
	// Set up the index page
	r.GET("/", func(c *gin.Context) {
		// Serve the emails we generated
		readfile, err := ioutil.ReadFile("news.json")
		if err != nil {
			log.Panic(err)
		}
		var decoded []news
		json.Unmarshal(readfile, &decoded)
		c.HTML(http.StatusOK, "index.html", gin.H{"domain": domain, "prettyname": prettyname, "news": decoded})
	})
	r.GET("/login", func(c *gin.Context) { c.HTML(http.StatusOK, "login.html", gin.H{"error": nil}) })
	r.GET("/register", func(c *gin.Context) { c.HTML(http.StatusOK, "register.html", nil) })
	r.GET("/inbox", func(c *gin.Context) {
		names := utils.Generatemail(15) // Generate 15 emails
		c.HTML(http.StatusOK, "mail.html", gin.H{"names": names})
	})
	r.POST("/login", func(c *gin.Context) {
		// Get the username and password
		username := c.PostForm("username")
		password := c.PostForm("password")
		hash := sha256.Sum256([]byte(password))
		encpw := hex.EncodeToString(hash[:])
		email := username + "@" + domain
		dbres := db.UUIDfromemail(email)
		if dbres == "" {
			c.HTML(http.StatusOK, "login.html", gin.H{"error": "Invalid login"})
		} else {

			user := sqlstuff.User{UUID: dbres, Name: username, Email: email, Password: encpw}
			if db.Get(user) == nil {
				c.HTML(http.StatusOK, "login.html", gin.H{"error": "Invalid login"})
			} else {
				c.Redirect(302, "../inbox")
			}
		}
	})
	// Same as above but for the register page
	r.POST("/register", func(c *gin.Context) {
		// Get the username and password
		username := c.PostForm("username")
		password := c.PostForm("password")
		hash := sha256.Sum256([]byte(password))
		encpw := hex.EncodeToString(hash[:])
		email := username + "@" + domain
		if utils.Validemail(email) == false {
			c.HTML(http.StatusOK, "register.html", gin.H{"error": "Invalid email/username"})
		}
		usr := sqlstuff.User{UUID: utils.CreateUUID(32), Name: username, Email: email, Password: encpw}
		res := db.InsertUser(usr)
		if res == true {
			c.Redirect(302, "/login")
		} else {
			c.HTML(http.StatusOK, "register.html", gin.H{"error": "Username already exists"})
		}
	})
	r.Use(gzip.Gzip(gzip.DefaultCompression)) // Gzip all the things
	r.Run()                                   // Run the server
}
