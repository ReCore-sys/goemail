package sqlstuff

import (
	"database/sql"
	"log"

	"github.com/ReCore-sys/gomail/libraries/utils"
)

// User is a struct to store the data for the users
type User struct {
	UUID     string
	Name     string
	Email    string
	Password string
}

// SQL is a struct to store the database connection
type SQL struct {
	Sqlcon    *sql.DB
	Tablename string
}

// Close closes the database connection
func (s *SQL) Close() error {
	return s.Sqlcon.Close()
}

// InsertUser adds a new user to the database
func (s *SQL) InsertUser(user User) bool {
	tx, err := s.Sqlcon.Begin()
	all := s.Getallfromcollumn("Email")
	if utils.ContainsString(all, user.Email) {
		log.Println("Email already exists")
		return false
	}
	if err != nil {
		log.Panic(err)
		return false
	}
	stmt, err := tx.Prepare("INSERT INTO " + s.Tablename + " VALUES(?, ?, ?, ?)")
	if err != nil {
		log.Panic(err)
		return false
	}
	_, err = stmt.Exec(user.UUID, user.Name, user.Email, user.Password)
	if err != nil {
		log.Panic(err)
		return false

	}
	tx.Commit()
	return true
}

// Get grab the needed data from the database for a specific email
func (s *SQL) Get(user User) []string {
	query := `SELECT * FROM ` + s.Tablename + ` WHERE email = '` + user.Email + "'"
	rows, err := s.Sqlcon.Query(query)
	if err != nil {
		log.Panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var UUID string
		var name string
		var email string
		var password string
		err = rows.Scan(&UUID, &name, &email, &password)
		if err != nil {
			log.Panic(err)
		}
		if email == email {
			user.Name = name
			user.Email = email
			user.Password = password
			return []string{user.Name, user.Email, user.Password}
		}
	}
	return nil
}

// Getcollumn finds and returns a specific collumn from the database
func (s *SQL) Getcollumn(user User, collumn string) string {
	query := `SELECT ` + collumn + ` FROM ` + s.Tablename + ` WHERE UUID = ` + user.UUID
	rows, err := s.Sqlcon.Query(query)
	if err != nil {
		println("Query issue")
		log.Panic(err)
	}
	defer rows.Close()
	var collumnname string
	for rows.Next() {
		err = rows.Scan(&collumnname)
		if err != nil {
			println("Scan issue")
			log.Panic(err)
		}
	}
	return collumnname

}

// Getallfromcollumn finds and returns all the data from a specific collumn from the database
func (s *SQL) Getallfromcollumn(collumn string) []string {
	query := `SELECT ` + collumn + ` FROM ` + s.Tablename
	rows, err := s.Sqlcon.Query(query)
	if err != nil {
		log.Panic(err)
	}
	defer rows.Close()
	var values []string
	for rows.Next() {
		var value string
		err = rows.Scan(&value)
		if err != nil {
			log.Panic(err)
		}
		values = append(values, value)
	}
	return values
}

// UUIDfromemail finds the UUID from an email
func (s *SQL) UUIDfromemail(email string) string {
	tx, err := s.Sqlcon.Begin()
	if err != nil {
		log.Panic(err)
	}
	q, err := tx.Prepare(`SELECT UUID FROM ` + s.Tablename + " WHERE email = '" + email + "'")
	if err != nil {
		log.Panic(err)
	}
	rows, err := q.Query()
	if err != nil {
		log.Panic(err)
	}
	var uuid string
	for rows.Next() {
		err = rows.Scan(&uuid)
		if err != nil {
			log.Panic(err)
		}
	}
	rows.Close()
	return uuid
}

// UserfromUUID does what it says on the tin
func (s *SQL) UserfromUUID(uuid string) User {
	tx, err := s.Sqlcon.Begin()
	if err != nil {
		log.Panic(err)
	}
	q, err := tx.Prepare(`SELECT * FROM ` + s.Tablename + " WHERE UUID = '" + uuid + "'")
	if err != nil {
		log.Panic(err)
	}
	rows, err := q.Query()
	if err != nil {
		log.Panic(err)
	}
	var user User
	for rows.Next() {
		err = rows.Scan(&user.UUID, &user.Name, &user.Email, &user.Password)
		if err != nil {
			log.Panic(err)
		}
	}
	rows.Close()
	return user
}
