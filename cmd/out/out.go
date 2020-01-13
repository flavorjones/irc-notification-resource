package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/flavorjones/irc-notification-resource/pkg/irc"
)

func exitWithError(error interface{}) {
	fmt.Fprintf(os.Stderr, "ERROR: irc-notification-resource/out: %v\n", error)
	os.Exit(1)
}

func main() {
	var (
		destPath string
		data     []byte
		message  string
		request  *irc.Request
		err      error
	)

	request, err = irc.ParseAndCheckRequest(os.Stdin)
	if err != nil {
		exitWithError(err)
	}

	if request.Params.Message != "" {
		message = request.Params.Message
	}

	if request.Params.MessageFile != "" {
		destPath = path.Join(os.Args[1], request.Params.MessageFile)
		data, err = ioutil.ReadFile(destPath)

		if err != nil {
			exitWithError(err)
		}

		message = string(data)
	}

	message = irc.ExpandMessage(message)

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
