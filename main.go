package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

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

func init() {
	if secret == "" {
		log.Fatalf("[fatal] invalid secret: %+q", secret)
	}
	if addr == "" {
		log.Fatalf("[fatal] invalid addr: %+q", addr)
	}
}

func main() {
	hook, _ := github.New(github.Options.Secret(secret))

	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		payload, err := hook.Parse(r, github.PushEvent, github.PullRequestEvent)
		if err != nil {
			log.Printf("[receiving] hook.Parse: err=%v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		switch p := payload.(type) {

		case github.PushPayload:
			// fmt.Printf("[receiving] push: %+v\n", push) // ** verbose
			fmt.Printf("[receiving] push: repo=%s, sender=%s\n", p.Repository.FullName, p.Sender.Login)

			prjKey := strings.ReplaceAll(p.Repository.FullName, "/", "_")
			prjKey = strings.ReplaceAll(prjKey, "-", "_")
			envKey := fmt.Sprintf("PRJ_PATH_%s", prjKey)
			prjPath := os.Getenv(envKey)

			if err := pull(prjPath); err != nil {
				log.Printf("[error] pull: path=%q, err=%v\n", prjPath, err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			fmt.Printf("[done] pull: path=%q\n", prjPath)

		case github.PullRequestPayload:
			// Do whatever you want from here...
			fmt.Printf("[receiving] pullRequest: repo=%s, sender=%s\n", p.Repository.FullName, p.Sender.Login)
		}
	})
	http.ListenAndServe(addr, nil)
}
