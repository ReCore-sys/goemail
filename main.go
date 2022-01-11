package main

import (
	"bufio"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"

	"github.com/Pallinder/go-randomdata"
	"github.com/gin-gonic/gin"
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

// Here we store a smaller bit of email
type email struct {
	Header  string `json:"header"`
	Content string `json:"content"`
}

func generatemail() []mails {

	// Load the mails struct
	names := []mails{}
	// Read the file of spam emails
	lines, _ := readLines("C:/Users/ReCor/Documents/OtherCode/goemail/spams_new.txt")
	// Shuffle it
	rand.Shuffle(len(lines), func(i, j int) { lines[i], lines[j] = lines[j], lines[i] })
	// Set up the struct to store all the possible emails
	possibles := []email{}
	// Loop through all the possible ones
	for _, line := range lines {
		// Split the line into the header and the content
		linecut := strings.Split(line, "|")
		// Trim whitespace
		for i, linecut2 := range linecut {
			linecut[i] = strings.TrimSpace(linecut2)
		}
		// Add it to the list and keep going
		possibles = append(possibles, email{linecut[0], linecut[1]})
	}
	// Generate 10 emails
	for i := 0; i < 10; i++ {
		// Get a random name
		pname := randomdata.FullName(randomdata.RandomGender)
		// Get a random email
		domain := randomdata.Email()
		// Get a random date
		date := randomdata.FullDate()
		// Pick a random header and content
		randomchoice := possibles[rand.Intn(len(possibles))]
		// Get all the possible images
		choiced := getdircontent("C:/Users/ReCor/Documents/OtherCode/goemail/static/img")
		// Pick a random image
		img := choiced[rand.Intn(len(choiced))]
		// Slap it all together
		data := mails{randomchoice.Header, randomchoice.Content, pname, domain, date, img}
		// Throw it into the slice and start again!
		names = append(names, data)
	}
	return names
}

// Set up a reader to read the file line by line.
func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

// Get a slice of all the files in a dir
func getdircontent(dir string) (content []string) {
	content = make([]string, 0)
	files, _ := ioutil.ReadDir(dir)

	for _, f := range files {
		content = append(content, f.Name())
	}
	return content
}

// Letsa go!
func main() {

	// Set up gin
	r := gin.Default()
	// Set up paths to the templates and the static files
	r.LoadHTMLGlob("templates/*.html")
	r.Static("/css", "static/css")
	r.Static("/js", "static/js")
	r.Static("/img", "static/img")
	// Set up the index page
	r.GET("/", func(c *gin.Context) {
		names := generatemail()
		// Server the page we just made
		c.HTML(http.StatusOK, "index.html", gin.H{"names": names})
	})
	r.GET("/login", func(c *gin.Context) { c.HTML(http.StatusOK, "login.html", nil) })

	r.Run()
}
