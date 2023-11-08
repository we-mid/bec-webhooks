package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-playground/webhooks/v6/github"
	_ "github.com/joho/godotenv/autoload"
)

const (
	path = "/webhooks"
)

var (
	secret = os.Getenv("SECRET")
	addr   = os.Getenv("ADDR")
)

func main() {
	hook, _ := github.New(github.Options.Secret(secret))

	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		payload, err := hook.Parse(r, github.PushEvent, github.PullRequestEvent)
		if err != nil {
			if err == github.ErrEventNotFound {
				// ok event wasn't one of the ones asked to be parsed
				fmt.Printf("** ErrEventNotFound: %+v\n", payload)
			}
		}
		switch payload.(type) {

		case github.PushPayload:
			push := payload.(github.PushPayload)
			// Do whatever you want from here...
			fmt.Printf("** push: %+v\n", push)

		case github.PullRequestPayload:
			pullRequest := payload.(github.PullRequestPayload)
			// Do whatever you want from here...
			fmt.Printf("** pullRequest: %+v\n", pullRequest)
		}
	})
	http.ListenAndServe(addr, nil)
}
