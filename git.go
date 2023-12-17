package main

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
)

func gitPullExt(dir string) (err error) {
	// Check if the git workspace is clean
	var output string
	output, err = runCmds(dir, []string{"git status --porcelain"})
	if err != nil {
		return
	}
	isClean := len(strings.TrimSpace(string(output))) == 0

	// Create a timestamp
	timestamp := time.Now().Format("2023-12-17_16:30:05.000")

	if !isClean {
		// 1. git stash -u with a timestamp as identity
		_, err = runCmds(dir, []string{
			fmt.Sprintf("git stash -u -m stash-%s", timestamp),
		})
		if err != nil {
			return
		}
	}
	finally := func() {
		if !isClean {
			_, e := runCmds(dir, []string{
				// 3. git stash pop with the last timestamp given
				fmt.Sprintf("git stash apply stash^{/stash-%s}", timestamp),
				// 4. if there is any conflicts, leave it for manually resolving
				"git add -A",
				"git reset",
			})
			if err == nil {
				err = e
			}
		}
	}
	defer finally()

	// 2. git pull --rebase
	_, err = runCmds(dir, []string{"git pull --rebase"})
	return
}

func runCmds(dir string, commands []string) (output string, err error) {
	for _, command := range commands {
		cmd := exec.Command("/bin/sh", "-c", command)
		cmd.Dir = dir
		var out []byte
		out, err = cmd.CombinedOutput()
		if err != nil {
			return
		}
		output += string(out)
	}
	return output, err
}
