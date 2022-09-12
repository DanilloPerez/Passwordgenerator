package main

import (
	"flag"
	"math/rand"
	"time"
)

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

	password := generatePass(lengte)
	println(password)
	println(specialCharacters, getallen)
}

func generatePass(lengte int) string {
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
	return password
}

func ConnectToDB() {

}

func CheckForExistingPass() {

}

func AddPass() {

}
