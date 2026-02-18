package main

import (
	"register-deployed-artifact-test/cmd"
	"log"
)

func main() {

	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
