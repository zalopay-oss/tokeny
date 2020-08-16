package main

import (
	"github.com/ltpquang/tokeny/pkg/totp"
	"log"
)

func main() {
	result, err := totp.Generate("")
	if err != nil {
		log.Panic(err)
	}
	log.Print(result)
}
