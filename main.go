package main

import (
	"log"
	"net/http"
	"os"

	"github.com/bortexel/stats-server/database"
)

const defaultBindAddr = ":8080"

var ConfiguredAuthorizationMiddleware func(ActionHandler) ActionHandler

func main() {
	uri, ok := os.LookupEnv("MONGO_CONNECTION_URI")
	if !ok {
		log.Fatalln("MONGO_CONNECTION_URI variable is empty")
		return
	}

	err := database.InitDatabase(uri)
	if err != nil {
		log.Println("Unable to init database:", err)
		return
	}

	bindAddr := os.Getenv("BIND_ADDR")
	if bindAddr == "" {
		bindAddr = defaultBindAddr
	}

	key, ok := os.LookupEnv("MUTATION_KEY")
	if !ok {
		log.Fatalln("MUTATION_KEY is not set")
	} else {
		ConfiguredAuthorizationMiddleware = Authorization(key)
	}

	log.Println("Starting HTTP server listener on", bindAddr)
	err = http.ListenAndServe(bindAddr, http.HandlerFunc(MainHandler))
	if err != nil {
		log.Fatalln("Error serving HTTP", err)
		return
	}
}
