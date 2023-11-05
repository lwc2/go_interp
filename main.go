package main

import (
	"fmt"
	"os"
	user2 "os/user"

	"go_interp/interp/repl"
)

func main() {
	user, err := user2.Current()

	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s! This is the Monkey programming language\n", user.Username)
	fmt.Printf("Feel free to type in commands\n")
	repl.Start(os.Stdin, os.Stdout)
}