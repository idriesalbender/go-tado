package main

import (
	"fmt"

	"github.com/idriesalbender/go-tado/tado"
)

func main() {
	// create a new tado client
	client := tado.NewClient()

	// get the authenticated user
	me, err := client.User.Get()
	if err != nil {
		panic(err)
	}

	// print a greeting
	fmt.Printf("Hello %s!\n", me.Name)
}
