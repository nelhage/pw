package main

import (
	"github.com/nelhage/pw"
	"log"
)

var config pw.Config = pw.Config{
	GPGKey:  "C808 7020 87F6 8CD8 C818  F239 DFC1 CF0D A816 9ACF",
	RootDir: "/home/nelhage/sec/pw",
}

func main() {
	log.Printf("GPG key: %s", config.GPGKey)
}
