package main

import (
	"context"
	"fmt"

	"github.com/idriesalbender/go-tado/tado"
)

func main() {
	// create a new tado client
	client := tado.NewClient()

	// get the authenticated user
	me, err := client.User.Get(context.Background())
	if err != nil {
		panic(err)
	}

	// print a greeting
	fmt.Printf("Hello %s!\n", me.Name)
}
