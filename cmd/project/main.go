package main

import (
	"log"

	"github.com/mnocard/go-project/internal/app"
)

func main() {
	log.Fatal(app.GetRouter())
}
