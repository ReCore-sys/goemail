package utils

import (
	"bufio"
	"io/ioutil"
	"log"
	"math/rand"
	"net/mail"
	"os"
	"strings"
	"time"

	"github.com/Pallinder/go-randomdata"
)

var filelocation, _ = os.Getwd()

// CheckFile is simple function to find if a file exists
func CheckFile(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

// Getdircontent gets a slice of all the files in a dir
func Getdircontent(dir string) (content []string) {
	content = make([]string, 0)
	files, _ := ioutil.ReadDir(dir)

	for _, f := range files {
		content = append(content, f.Name())
	}
	return content
}

// Here we store a smaller bit of email
type email struct {
	Header  string `json:"header"`
	Content string `json:"content"`
}

// Mails is a struct to store the data for the spam emails
type Mails struct {
	Message  string
	Content  string
	Name     string
	Domain   string
	Datetime string
	Img      string
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

// Generatemail Creates 10 randomised spam emails
func Generatemail(count int) []Mails {
	rand.Seed(time.Now().UnixNano())

	// Load the Mails struct
	names := []Mails{}
	// Read the file of spam eMails
	lines, err := readLines(filelocation + "/spams.txt")
	if err != nil {
		log.Panic(err)
	}
	// Shuffle it
	rand.Shuffle(len(lines), func(i, j int) { lines[i], lines[j] = lines[j], lines[i] })
	// Set up the struct to store all the possible eMails
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
	// Generate `count` eMails
	for i := 0; i < count; i++ {
		// Get a random name
		pname := randomdata.FullName(randomdata.RandomGender)
		// Get a random email
		domain := randomdata.Email()
		// Get a random date
		date := randomdata.FullDate()
		// Pick a random header and content
		var randomchoice email
		if len(possibles) > 0 {
			randomchoice = possibles[rand.Intn(len(possibles))]
		} else {
			randomchoice = email{"No content found", "No content found"}
		}
		// Get all the possible images
		choiced := Getdircontent(filelocation + "/static/img")
		// Pick a random image
		img := choiced[rand.Intn(len(choiced))]
		// Slap it all together
		data := Mails{randomchoice.Header, randomchoice.Content, pname, domain, date, img}
		// Throw it into the slice and start again!
		names = append(names, data)
	}
	return names
}

// RandomChoice returns a random selection from the given slice
func RandomChoice(inp []interface{}) interface{} {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(len(inp))
}

// CreateUUID creates a random UUID
func CreateUUID(length int) string {
	rand.Seed(time.Now().UnixNano())
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	var uuid string
	for i := 0; i < length; i++ {
		uuid += string(letters[rand.Intn(len(letters))])
	}
	return uuid
}

// ContainsInt returns true if the given int is in the given slice
func ContainsInt(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// ContainsString returns true if the given string is in the given slice
func ContainsString(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// Validemail checks if the email provided is valid
func Validemail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
