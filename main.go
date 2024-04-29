package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/nishokbanand/interpreter/repl"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Welcome to the language of the GODS, Mr.%v\n", user.Username)
	repl.Start(os.Stdin, os.Stdout)
}
