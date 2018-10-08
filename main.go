package main

import (
	"flag"
	"log"
	"os"
)

var (
	googleCredentialsFilePath = flag.String("google-credentials", "./credentials.json", "path to credentials file (default './credentials.json')")
	hangoutsServerPort        = flag.Int("hangouts-server-port", 8080, "port to run the Hangouts HTTP server on (default 8080)")
)

func main() {
	flag.Parse()

	// set google credentials via environment variable using file
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", *googleCredentialsFilePath)

	// initialize chat service and register plugins
	cs, err := NewChatService(*hangoutsServerPort)
	if err != nil {
		log.Fatal(err)
	}

	cs.RegisterPlugin(new(HelloWorldPlugin))

	// start chat service
	log.Fatal(cs.Serve())
}
