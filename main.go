package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"

	_ "github.com/lib/pq"
)

const (
	dbname = "mijndb"
	dbuser = "danillo"
	dbpass = "danillo"
)

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
	println(specialCharacters, getallen)
	ConnectToDB()
	CreateTable()
	password := GeneratePass(lengte)
	println(password)

}

func GeneratePass(lengte int) string {
	characters := "abcdefghijklmnopqrstuvwxyz"

	var password string

	if specialCharacters == true {
		characters += "+-()=!#$%^&"
	}
	if getallen == true {
		characters += "1234567890"
	}

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < lengte; i++ {
		password += string(characters[rand.Intn(len(characters))])

	}
	//password = "tmpfgrxu"
	//recursion
	exists, err := CheckForExistingPass(password)
	if err != nil {
		log.Fatal(err)
	}
	if exists {
		return GeneratePass(lengte)
	}

	err = AddPass(password)
	if err != nil {
		log.Fatal(err)
	}
	return password
}

func ConnectToDB() {
	db, err := sql.Open("postgres", "dbname="+dbname+" user="+dbuser+" password="+dbpass+" sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	database = db
	if err := CreateTable(); err != nil {
		log.Fatal(err)
	}
}

func CreateTable() error {
	query := `CREATE TABLE IF NOT EXISTS passwords (
password varchar(255)
	)`
	_, err := database.Exec(query)
	return err
}

func CheckForExistingPass(password string) (bool, error) {
	var exists bool
	query := fmt.Sprintf(`SELECT EXISTS(SELECT 1 FROM passwords WHERE password = '%s')`, password)
	err := database.QueryRow(query).Scan(&exists)
	println(exists)
	return exists, err
}

func AddPass(password string) error {
	query := fmt.Sprintf(`INSERT INTO passwords (password)VALUES ('%s')`, password)
	_, err := database.Exec(query)
	return err
}
