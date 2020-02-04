package main

import (
	"fmt"
	"monkey/repl"
	"os"
	"os/user"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Howdy %s, Welcome to Monkey programming language \n", user.Username)
	fmt.Printf("Feel free to type coments for Lexing! ")
	repl.Start(os.Stdin, os.Stdout)
}
