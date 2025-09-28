package main

import (
	"fmt"

	"github.com/sidis405/gator/internal/config"
)

func main() {
	configuration, _ := config.Read()
	err := configuration.SetUser("sid")

	if err != nil {
		fmt.Printf("%q", err)
		return
	}

	configuration, err = config.Read()
	if err != nil {
		fmt.Printf("%q", err)
		return
	}
	fmt.Println(configuration)
}
