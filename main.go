package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/thingsme/thingscript/repl"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s! thingscript is ready\n", user.Username)
	repl.Start(os.Stdin, os.Stdout)
}
