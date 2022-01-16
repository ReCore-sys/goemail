package utils

import (
	"bufio"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"

	"github.com/Pallinder/go-randomdata"
)

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
func Generatemail() []Mails {

	// Load the Mails struct
	names := []Mails{}
	// Read the file of spam eMails
	lines, _ := readLines("C:/Users/ReCor/Documents/OtherCode/gomail/spams_new.txt")
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
	// Generate 10 eMails
	for i := 0; i < 10; i++ {
		// Get a random name
		pname := randomdata.FullName(randomdata.RandomGender)
		// Get a random email
		domain := randomdata.Email()
		// Get a random date
		date := randomdata.FullDate()
		// Pick a random header and content
		var randomchoice email
		fmt.Println(len(possibles))
		if len(possibles) > 0 {
			randomchoice = possibles[rand.Intn(len(possibles))]
		} else {
			randomchoice = email{"No content found", "No content found"}
		}
		// Get all the possible images
		choiced := Getdircontent("C:/Users/ReCor/Documents/OtherCode/gomail/static/img")
		// Pick a random image
		img := choiced[rand.Intn(len(choiced))]
		// Slap it all together
		data := Mails{randomchoice.Header, randomchoice.Content, pname, domain, date, img}
		// Throw it into the slice and start again!
		names = append(names, data)
	}
	return names
}

// User is a struct to store the data for the users
type User struct {
	Name     string
	email    string
	password string
}

// SQL is a struct to store the database connection
type SQL struct {
	sqlcon    sql.DB
	tableName string
}

// Close closes the database connection
func (s *SQL) Close() error {
	return s.sqlcon.Close()
}

// InsertUser adds a new user to the database
func (s *SQL) InsertUser(user User) error {
	tx, err := s.sqlcon.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("INSERT INTO" + s.tableName + "VALUES(?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(user.Name, user.email, user.password)
	if err != nil {
		log.Fatal(err)

	}
	return tx.Commit()
}

// Select grab the needed data from the database for a specific email
func (s *SQL) Select(user User, email string) []string {
	query := `SELECT * FROM ` + s.tableName + ` WHERE email = ` + email
	rows, err := s.sqlcon.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		var email string
		var password string
		err = rows.Scan(&name, &email, &password)
		if err != nil {
			log.Fatal(err)
		}
		if email == email {
			user.Name = name
			user.email = email
			user.password = password
			return []string{user.Name, user.email, user.password}
		}
	}
	return nil
}
