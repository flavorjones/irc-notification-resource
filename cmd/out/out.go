package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/flavorjones/irc-notification-resource/pkg/irc"
)

func exitWithError(error interface{}) {
	fmt.Fprintf(os.Stderr, "ERROR: irc-notification-resource/out: %v\n", error)
	os.Exit(1)
}

func main() {
	request, err := irc.ParseAndCheckRequest(os.Stdin)
	if err != nil {
		exitWithError(err)
	}

	message := irc.ExpandMessage(request)

	if !request.Params.DryRun {
		irc.SendMessage(request, message)
	}

	response := irc.BuildResponse(request, message)

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(response)
	if err != nil {
		exitWithError(err)
	}
}
