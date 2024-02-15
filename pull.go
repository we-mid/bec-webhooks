package main

import (
	"fmt"
	"os/exec"
	"time"
)

func pull(dir string) error {
	timestamp := time.Now().Format("2023-12-17_16:30:05.000")
	sh := fmt.Sprintf(`
		cd %q
		not_clean=$(if [ -n "$(git status --porcelain)" ]; then echo true; else echo false; fi)
		dt=%q
		if [ $not_clean = true ]
		then
			git stash push -u -m "$dt"
		fi
		git pull --rebase
		if [ $not_clean = true ]
		then
			git stash pop
		fi
		git add -A
		git reset
	`, dir, timestamp)
	cmd := exec.Command("/bin/sh", "-c", sh)
	return cmd.Run()
}
