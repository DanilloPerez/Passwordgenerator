package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log"
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
	err := ConnectToDB()
	if err != nil {
		log.Fatal(err)
	}

	// generate pass will create a unique password for us and store it in the database
	password, err := GeneratePass(lengte)
	if err != nil {
		log.Fatal(err)
	}
	println(password)
}

func (cfg *config) GetConfig() *config {
	conf, err := os.ReadFile("conf.yml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}

	err = yaml.Unmarshal([]byte(conf), &cfg)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	return cfg
}

func GeneratePass(lengte int) (string, error) {
	if lengte <= 0 {
		return "", errors.New("Length of password can't be 0")
	}

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
	exists, err := CheckForExistingPass(password)
	if err != nil {
		return "", err
	}
	if exists {
		return GeneratePass(lengte)
	}

	// if the password does not exist yet in the database we add it
	err = AddPass(password)
	if err != nil {
		return "", err
	}
	return password, nil
}

func ConnectToDB() error {
	var c config
	c.GetConfig()

	// checking database credentials for empty values
	if IsEmpty(c.Dbname) || IsEmpty(c.Dbpass) || IsEmpty(c.Dbuser) {
		return errors.New("Database credentials were not supplied correctly")
	}
	// open connection to database
	db, err := sql.Open("postgres", "dbname="+c.Dbname+" user="+c.Dbuser+" password="+c.Dbpass+" sslmode=disable")
	if err != nil {
		return errors.New("Connection to the database could not be established")
	}
	database = db

	// create required tables if they do not yet exist
	if err := CreateTable(); err != nil {
		return err
	}
	return nil
}

func CreateTable() error {
	query := `CREATE TABLE IF NOT EXISTS passwords (password varchar(255))`
	_, err := database.Exec(query)
	if err != nil {
		return errors.New("An error occured while creating your table")
	}
	return nil
}

func CheckForExistingPass(password string) (bool, error) {
	var exists bool
	query := fmt.Sprintf(`SELECT EXISTS(SELECT 1 FROM passwords WHERE password = '%s')`, password)

	// execute query and bind assign boolean return value to the "exists" variable
	err := database.QueryRow(query).Scan(&exists)
	if err != nil {
		return exists, errors.New("An error occured while saving your password to the database")
	}
	return exists, nil
}

func AddPass(password string) error {
	query := fmt.Sprintf(`INSERT INTO passwords (password)VALUES ('%s')`, password)
	_, err := database.Exec(query)
	if err != nil {
		return errors.New("An error occured while saving your password to the database")
	}
	return nil
}

func IsEmpty(textValue string) bool {
	// trim leading and trailing whitespaces from string and check length to verify it contains a value
	return (len(strings.TrimSpace(textValue)) == 0)
}
