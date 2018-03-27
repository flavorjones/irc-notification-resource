package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/Elemental-IRCd/irc"
)

//
//  structs for reading and writing our json
//
type Source struct {
	Server   string `json:"server"`
	Port     int    `json:"port"`
	Channel  string `json:"channel"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type Params struct {
	Message string `json:"message"`
	DryRun  bool   `json:"dry_run"` // undocumented
}

type Request struct {
	Source Source `json:"source"`
	Params Params `json:"params"`
}

type Version struct {
	Ref string `json:"ref"`
}

type Metadatum struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

type Response struct {
	Version  Version     `json:"version"`
	Metadata []Metadatum `json:"metadata"`
}

func exitWithError(error interface{}) {
	fmt.Fprintf(os.Stderr, "ERROR: irc-notification-resource/out: %v\n", error)
	os.Exit(1)
}

const buildUrlTemplate = "${ATC_EXTERNAL_URL}/teams/${BUILD_TEAM_NAME}/pipelines/${BUILD_PIPELINE_NAME}/jobs/${BUILD_JOB_NAME}/builds/${BUILD_NAME}"

func ParseAndCheckRequest(reader io.Reader) (*Request, error) {
	request := Request{}

	err := json.NewDecoder(reader).Decode(&request)
	if err != nil {
		return &request, err
	}

	if request.Source.Server == "" {
		return &request, errors.New("No server was provided")
	}
	if request.Source.Port == 0 {
		return &request, errors.New("No port was provided")
	}

	if request.Source.Channel == "" {
		return &request, errors.New("No channel was provided")
	}
	if request.Source.User == "" {
		return &request, errors.New("No user was provided")
	}
	if request.Source.Password == "" {
		return &request, errors.New("No password was provided")
	}
	if request.Params.Message == "" {
		return &request, errors.New("No message was provided")
	}

	return &request, nil
}

func ExpandMessage(request *Request) string {
	os.Setenv("BUILD_URL", os.ExpandEnv(buildUrlTemplate))
	return os.ExpandEnv(request.Params.Message)
}

func BuildResponse(request *Request, message string) *Response {
	// omit password for reasons that are hopefully obvious
	response := Response{}
	response.Metadata = append(response.Metadata, Metadatum{"host", connString(request)})
	response.Metadata = append(response.Metadata, Metadatum{"channel", request.Source.Channel})
	response.Metadata = append(response.Metadata, Metadatum{"user", request.Source.User})
	response.Metadata = append(response.Metadata, Metadatum{"message", message})
	response.Metadata = append(response.Metadata, Metadatum{"dry_run", request.Params.DryRun})
	response.Version.Ref = "none"
	return &response
}

func SendMessage(request *Request, message string) error {
	conn := irc.New(request.Source.User, request.Source.User)

	conn.UseTLS = true
	conn.Log = log.New(ioutil.Discard, "", 0) // be completely silent
	conn.Password = request.Source.Password

	conn.AddCallback("001", func(*irc.Event) {
		conn.Privmsg(request.Source.Channel, message)
		conn.Quit()
	})

	err := conn.Connect(connString(request))
	if err != nil {
		return err
	}

	conn.Loop()

	return nil
}

func connString(request *Request) string {
	return fmt.Sprintf("%s:%d", request.Source.Server, request.Source.Port)
}

func main() {
	request, err := ParseAndCheckRequest(os.Stdin)
	if err != nil {
		exitWithError(err)
	}

	message := ExpandMessage(request)

	if request.Params.DryRun {
		fmt.Fprintf(os.Stderr, "NOTICE: dry run: not sending: `%s`\n", message)
	} else {
		SendMessage(request, message)
	}

	response := BuildResponse(request, message)

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(response)
	if err != nil {
		exitWithError(err)
	}
}
