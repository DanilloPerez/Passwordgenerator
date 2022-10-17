# Passwordgenerator
A generator for.... you guessed it, a password. 

How does it work ? 

The Go programme is booted with the terminal-command 'go run main.go'
After the command you can add parameters for the password you want generated. 
You have three options. 
The length which can be added with the '-l' command.
The addition of numbers with can be added with the '-g=true' command
and lastly the addition of special characters which can be added with the '-t=true' command. 

For example:
We want a password that may contain both special characters and numbers and is 13 characters long.
The command in this case should be: 

go run main.go -l 13 -g=true -t=true

----------------------------------------------------------------------------------------------------------------------------------------

How do you adjust database credentials ?

This programme makes use of a database for storing a password and checking if the password already exists.
If need arrises for use of a different database, you can adjust the credentials in the config-file. 
This works with the following 'commands'.
The current database works with a Dbname, Dbuser and Dbpass. 
you can fill in the credentials between the quote's

for example:

Dbname: "Database" 
Dbuser: "Vlad P."
Dbpass: "SpecialOperation"




