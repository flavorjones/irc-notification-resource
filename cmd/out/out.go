package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/Elemental-IRCd/irc"
)

//
//  structs for reading and writing our json
//
type source struct {
	Server   string `json:"server"`
	Port     int    `json:"port"`
	Channel  string `json:"channel"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type params struct {
	Message string `json:"message"`
	DryRun  bool   `json:"dry_run"` // undocumented
}

type request struct {
	Source source `json:"source"`
	Params params `json:"params"`
}

type version struct {
	Ref string `json:"ref"`
}

type metadatum struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

type response struct {
	Version  version     `json:"version"`
	Metadata []metadatum `json:"metadata"`
}

//
//  utility functions
//
func exitWithError(error interface{}) {
	fmt.Fprintf(os.Stderr, "ERROR: irc-notification-resource/out: %v\n", error)
	os.Exit(1)
}

//
//  constants
//
const buildUrlTemplate = "${ATC_EXTERNAL_URL}/teams/${BUILD_TEAM_NAME}/pipelines/${BUILD_PIPELINE_NAME}/jobs/${BUILD_JOB_NAME}/builds/${BUILD_NAME}"

func main() {
	request := request{}
	response := response{}

	err := json.NewDecoder(os.Stdin).Decode(&request)
	if err != nil {
		exitWithError(err)
	}

	if request.Source.Server == "" {
		exitWithError("No server was provided")
	}
	if request.Source.Port == 0 {
		exitWithError("No port was provided")
	}
	connString := fmt.Sprintf("%s:%d", request.Source.Server, request.Source.Port)

	if request.Source.Channel == "" {
		exitWithError("No channel was provided")
	}
	if request.Source.User == "" {
		exitWithError("No user was provided")
	}
	if request.Source.Password == "" {
		exitWithError("No password was provided")
	}
	if request.Params.Message == "" {
		exitWithError("No message was provided")
	}

	// set up an environment variable for the build url
	os.Setenv("BUILD_URL", os.ExpandEnv(buildUrlTemplate))

	// expand any environment variables in the message text
	message := os.ExpandEnv(request.Params.Message)

	if request.Params.DryRun {
		fmt.Fprintf(os.Stderr, "dry run: not sending '%s'\n", message)
	} else {
		conn := irc.New(request.Source.User, request.Source.User)

		conn.UseTLS = true
		conn.Log = log.New(ioutil.Discard, "", 0) // be completely silent
		conn.Password = request.Source.Password

		conn.AddCallback("001", func(*irc.Event) {
			conn.Privmsg(request.Source.Channel, message)
			conn.Quit()
		})

		err = conn.Connect(connString)
		if err != nil {
			exitWithError(err)
		}

		conn.Loop()
	}

	// build response struct
	// omit password for reasons that are hopefully obvious
	response.Metadata = append(response.Metadata, metadatum{"host", connString})
	response.Metadata = append(response.Metadata, metadatum{"channel", request.Source.Channel})
	response.Metadata = append(response.Metadata, metadatum{"user", request.Source.User})
	response.Metadata = append(response.Metadata, metadatum{"message", message})
	response.Version.Ref = "none"

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(response)
	if err != nil {
		exitWithError(err)
	}
}
