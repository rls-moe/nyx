package main

import (
	"fmt"
	"go.rls.moe/nyx/config"
	"go.rls.moe/nyx/http"
)

func main() {
	c, err := config.Load()
	if err != nil {
		fmt.Printf("Could not read configuration: %s\n", err)
		return
	}

	fmt.Println("Starting Server")
	http.Start(c)
}
