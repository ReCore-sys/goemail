package main

import (
	"database/sql"

	"io/ioutil"

	"net/http"
	"os"

	"log"

	gomail "github.com/ReCore-sys/gomail/libraries"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3" // We null this cos idfk
)

// Set up a struct to store info about a email
type mails struct {
	Message  string
	Content  string
	Name     string
	Domain   string
	Datetime string
	Img      string
}

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

// Letsa go!
func main() {
	// Create the database
	if gomail.CheckFile("data.db") == false {
		os.Create(filelocation + "/data.db")
	}
	db, err := sql.Open("sqlite3", filelocation+"/data.db")
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}
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
		names := gomail.Generatemail()
		// Server the page we just made
		c.HTML(http.StatusOK, "index.html", gin.H{"names": names})
	})
	r.GET("/login", func(c *gin.Context) { c.HTML(http.StatusOK, "login.html", nil) })
	r.GET("/register", func(c *gin.Context) { c.HTML(http.StatusOK, "register.html", nil) })
	r.POST("/login", func(c *gin.Context) {
		// Get the username and password
		username := c.PostForm("username")
		password := c.PostForm("password")
		print(username, password)
		c.Redirect(302, "/")
	})
	// Same as above but for the register page
	r.POST("/register", func(c *gin.Context) {
		// Get the username and password
		username := c.PostForm("username")
		password := c.PostForm("password")
		println("Got this far")

		println(username, password)
		c.Redirect(302, "/login")
	})
	r.Use(gzip.Gzip(gzip.DefaultCompression)) // Gzip all the things
	r.Run()                                   // Run the server
}
