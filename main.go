package main

import (
	"database/sql"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"gopkg.in/yaml.v3"
)

type config struct {
	Dbname string `yaml:"Dbname"`
	Dbuser string `yaml:"Dbuser"`
	Dbpass string `yaml:"Dbpass"`
}

var database *sql.DB

var lengte int
var getallen bool
var specialCharacters bool

func init() {
	flag.IntVar(&lengte, "l", 8, "ingevoerde lengte")
	flag.BoolVar(&getallen, "g", false, "ingevoerde getallen")
	flag.BoolVar(&specialCharacters, "t", false, "ingevoerde tekens")
	flag.Parse()
}

func main() {
	// connect to database and create required tables
	ConnectToDB()

	// create required tables if they do not yet exist
	CreateTable()

	// generate password will create a unique password for us and store it in the database
	GeneratePass(lengte)
}

func (cfg *config) GetConfig() {
	yamlFileName := "conf.yml"

	conf, err := os.ReadFile(yamlFileName)
	if err != nil {
		panic("Failed to read configuration from " + yamlFileName)
	}

	err = yaml.Unmarshal([]byte(conf), &cfg)
	if err != nil {
		panic("Could not bind data to Config struct")
	}
}

func GeneratePass(lengte int) {
	var password string

	// build string of to-be-used characters based on the user supplied flags (-g & -t)
	characters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	if specialCharacters == true {
		characters += "+-()=!#$%^&"
	}
	if getallen == true {
		characters += "1234567890"
	}

	rand.Seed(time.Now().UnixNano())
	// for however many characters the user requested with the -l flag:
	// we take a random character from the characters string and add it to the password string
	for i := 0; i < lengte; i++ {
		password += string(characters[rand.Intn(len(characters))])
	}

	// for DEV only
	//password = "tmpfgrxu"

	// recursion
	exists := CheckForExistingPass(password)
	if exists {
		GeneratePass(lengte)
	}

	// if the password does not exist yet in the database we add it
	AddPass(password)
	println("Password created: " + password)
}

func ConnectToDB() {
	var c config
	c.GetConfig()

	// checking database credentials for empty values
	if IsEmpty(c.Dbname) || IsEmpty(c.Dbpass) || IsEmpty(c.Dbuser) {
		panic("Database credentials were not supplied correctly")
	}

	// open connection to database
	db, err := sql.Open("postgres", "dbname="+c.Dbname+" user="+c.Dbuser+" password="+c.Dbpass+" sslmode=disable")
	if err != nil {
		panic("Connection to the database could not be established")
	}
	database = db
}

func CreateTable() {
	query := `CREATE TABLE IF NOT EXISTS passwords (password varchar(255))`
	_, err := database.Exec(query)
	if err != nil {
		panic("An error occured while creating your table")
	}
}

func CheckForExistingPass(password string) bool {
	var exists bool
	query := fmt.Sprintf(`SELECT EXISTS(SELECT 1 FROM passwords WHERE password = '%s')`, password)

	// execute query and bind assign boolean return value to the "exists" variable
	err := database.QueryRow(query).Scan(&exists)
	if err != nil {
		panic("An error occured while saving your password to the database")
	}
	return exists
}

func AddPass(password string) {
	query := fmt.Sprintf(`INSERT INTO passwords (password)VALUES ('%s')`, password)
	_, err := database.Exec(query)
	if err != nil {
		panic("An error occured while saving your password to the database")
	}
}

func IsEmpty(textValue string) bool {
	// trim leading and trailing whitespaces from string and check length to verify it contains a value
	return (len(strings.TrimSpace(textValue)) == 0)
}
