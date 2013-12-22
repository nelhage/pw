package main

import (
	"github.com/nelhage/pw"
	"log"
)

func main() {
	conf := pw.LoadConfig()
	log.Printf("GPG key: %s", conf.GPGKey)
}
