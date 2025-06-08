package main

import (
	"log"

	"github.com/mnocard/go-project/internal/app"
)

//	@title			Project Swagger
//	@version		0.1
//	@license.name	Apache 2.0

//	@host		localhost:8080
//	@BasePath	/
func main() {
	log.Fatal(app.RunRouter())
}
