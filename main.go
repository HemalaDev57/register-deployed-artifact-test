package main

import (
	"gha-register-deployed-artifact/cmd"
	"log"
)

func main() {

	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
