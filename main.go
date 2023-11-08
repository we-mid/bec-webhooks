package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
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
			if err == github.ErrEventNotFound {
				// ok event wasn't one of the ones asked to be parsed
				fmt.Printf("[receiving] ErrEventNotFound: %v\n", err)
			}
		}
		switch payload.(type) {

		case github.PushPayload:
			p := payload.(github.PushPayload)
			// Do whatever you want from here...
			// fmt.Printf("[receiving] push: %+v\n", push) // ** verbose
			fmt.Printf("[receiving] push: repo=%s, sender=%s\n", p.Repository.FullName, p.Sender.Login)

			prjKey := strings.ReplaceAll(p.Repository.FullName, "/", "_")
			prjKey = strings.ReplaceAll(prjKey, "-", "_")
			envKey := fmt.Sprintf("PRJ_PATH_%s", prjKey)
			prjPath := os.Getenv(envKey)
			cmds := fmt.Sprintf("cd %q && (git stash || true) && git pull --rebase && (git stash pop || true)", prjPath)

			cmd := exec.Command("/bin/sh", "-c", cmds)
			if err := cmd.Run(); err != nil {
				log.Printf("[error] git pull: path=%q, err=%v\n", prjPath, err)
				return
			}
			fmt.Printf("[done] git pull: path=%q\n", prjPath)

		case github.PullRequestPayload:
			p := payload.(github.PullRequestPayload)
			// Do whatever you want from here...
			// fmt.Printf("[receiving] pullRequest: %+v\n", pullRequest) // verbose
			fmt.Printf("[receiving] pullRequest: repo=%s, sender=%s\n", p.Repository.FullName, p.Sender.Login)
		}
	})
	http.ListenAndServe(addr, nil)
}
