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
			log.Printf("[receiving] hook.Parse: err=%v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		switch payload.(type) {

		case github.PushPayload:
			p := payload.(github.PushPayload)
			// fmt.Printf("[receiving] push: %+v\n", push) // ** verbose
			fmt.Printf("[receiving] push: repo=%s, sender=%s\n", p.Repository.FullName, p.Sender.Login)

			prjKey := strings.ReplaceAll(p.Repository.FullName, "/", "_")
			prjKey = strings.ReplaceAll(prjKey, "-", "_")
			envKey := fmt.Sprintf("PRJ_PATH_%s", prjKey)
			prjPath := os.Getenv(envKey)

			// todo: stash and pop by specific name
			// todo: detect if ws is clean
			cmds := []string{
				fmt.Sprintf("cd %q", prjPath),
				"git stash",
				"git pull --rebase",
				"(git stash pop || true)",
				// if there is any conflicts, leave it for manually resolving
				"git add .",
				"git reset .",
			}
			sh := strings.Join(cmds, " && ")
			cmd := exec.Command("/bin/sh", "-c", sh)
			if err := cmd.Run(); err != nil {
				log.Printf("[error] git pull: path=%q, err=%v\n", prjPath, err)
				w.WriteHeader(http.StatusInternalServerError)
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
