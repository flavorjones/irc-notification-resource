package main_test

import (
	"bytes"
	"os"

	. "github.com/flavorjones/irc-notification-resource/cmd/out"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

// {"source": {"server": "chat.freenode.net", "port": 7070, "channel": "#random", "user": "randobot1337", "password": "secretsecret"}}

var _ = Describe("Out", func() {
	Describe("request", func() {
		Context("when valid", func() {
			Describe("source", func() {
				It("returns correct values", func() {
					request, error := ParseAndCheckRequest(bytes.NewBufferString(`{"source": {"server": "chat.freenode.net", "port": 7070, "channel": "#random", "user": "randobot1337", "password": "secretsecret"}, "params": {"message": "foo"}}`))
					Expect(error).To(BeNil())
					Expect(request.Source).To(MatchAllFields(Fields{
						"Server":   Equal("chat.freenode.net"),
						"Port":     Equal(7070),
						"Channel":  Equal("#random"),
						"User":     Equal("randobot1337"),
						"Password": Equal("secretsecret"),
					}))
				})
			})

			Describe("params", func() {
				Describe("`message`", func() {
					It("returns correct value", func() {
						request, error := ParseAndCheckRequest(bytes.NewBufferString(`{"source": {"server": "chat.freenode.net", "port": 7070, "channel": "#random", "user": "randobot1337", "password": "secretsecret"}, "params": {"message": "foo $BUILD_ID"}}`))
						Expect(error).To(BeNil())
						Expect(request.Params).To(MatchFields(IgnoreExtras, Fields{
							"Message": Equal("foo $BUILD_ID"),
						}))
					})
				})

				Describe("`dry_run`", func() {
					It("defaults to false", func() {
						request, error := ParseAndCheckRequest(bytes.NewBufferString(`{"source": {"server": "chat.freenode.net", "port": 7070, "channel": "#random", "user": "randobot1337", "password": "secretsecret"}, "params": {"message": "foo $BUILD_ID"}}`))
						Expect(error).To(BeNil())
						Expect(request.Params).To(MatchFields(IgnoreExtras, Fields{"DryRun": BeFalse()}))
					})

					It("is settable to true", func() {
						request, error := ParseAndCheckRequest(bytes.NewBufferString(`{"source": {"server": "chat.freenode.net", "port": 7070, "channel": "#random", "user": "randobot1337", "password": "secretsecret"}, "params": {"message": "foo $BUILD_ID", "dry_run": true}}`))
						Expect(error).To(BeNil())
						Expect(request.Params).To(MatchFields(IgnoreExtras, Fields{"DryRun": BeTrue()}))
					})

					It("is settable to false", func() {
						request, error := ParseAndCheckRequest(bytes.NewBufferString(`{"source": {"server": "chat.freenode.net", "port": 7070, "channel": "#random", "user": "randobot1337", "password": "secretsecret"}, "params": {"message": "foo $BUILD_ID", "dry_run": false}}`))
						Expect(error).To(BeNil())
						Expect(request.Params).To(MatchFields(IgnoreExtras, Fields{"DryRun": BeFalse()}))
					})
				})
			})
		})

		Describe("required source property", func() {
			It("`server`", func() {
				_, error := ParseAndCheckRequest(bytes.NewBufferString(`{"source": {"port": 7070, "channel": "#random", "user": "randobot1337", "password": "secretsecret"}}`))
				Expect(error.Error()).To(MatchRegexp(`No server was provided`))
			})

			It("`port`", func() {
				_, error := ParseAndCheckRequest(bytes.NewBufferString(`{"source": {"server": "chat.freenode.net", "channel": "#random", "user": "randobot1337", "password": "secretsecret"}}`))
				Expect(error.Error()).To(MatchRegexp(`No port was provided`))
			})

			It("`channel`", func() {
				_, error := ParseAndCheckRequest(bytes.NewBufferString(`{"source": {"server": "chat.freenode.net", "port": 7070, "user": "randobot1337", "password": "secretsecret"}}`))
				Expect(error.Error()).To(MatchRegexp(`No channel was provided`))
			})

			It("`user`", func() {
				_, error := ParseAndCheckRequest(bytes.NewBufferString(`{"source": {"server": "chat.freenode.net", "port": 7070, "channel": "#random", "password": "secretsecret"}}`))
				Expect(error.Error()).To(MatchRegexp(`No user was provided`))
			})

			It("`password`", func() {
				_, error := ParseAndCheckRequest(bytes.NewBufferString(`{"source": {"server": "chat.freenode.net", "port": 7070, "channel": "#random", "user": "randobot1337"}}`))
				Expect(error.Error()).To(MatchRegexp(`No password was provided`))
			})
		})

		Describe("required params property", func() {
			It("`message`", func() {
				_, error := ParseAndCheckRequest(bytes.NewBufferString(`{"source": {"server": "chat.freenode.net", "port": 7070, "channel": "#random", "user": "randobot1337", "password": "secretsecret"}}`))
				Expect(error.Error()).To(MatchRegexp(`No message was provided`))
			})
		})
	})

	Describe("metadata expansion", func() {
		var request Request

		BeforeEach(func() {
			request = Request{
				Source: Source{
					Server:   "chat.freenode.net",
					Port:     7070,
					Channel:  "#random",
					User:     "randobot1337",
					Password: "secretsecret",
				},
				Params: Params{DryRun: true},
			}

			os.Setenv("BUILD_ID", "id-123")
			os.Setenv("BUILD_NAME", "name-asdf")
			os.Setenv("BUILD_JOB_NAME", "job-name-asdf")
			os.Setenv("BUILD_PIPELINE_NAME", "pipeline-name-asdf")
			os.Setenv("BUILD_TEAM_NAME", "team-name-asdf")
			os.Setenv("ATC_EXTERNAL_URL", "https://ci.example.com")
		})

		It("expands environment variables", func() {
			request.Params.Message = ">> $BUILD_ID <<"
			message := ExpandMessage(&request)
			Expect(message).To(Equal(">> id-123 <<"))
		})

		It("expands BUILD_URL pseudo-metadata", func() {
			request.Params.Message = ">> $BUILD_URL <<"
			message := ExpandMessage(&request)
			Expect(message).To(Equal(">> https://ci.example.com/teams/team-name-asdf/pipelines/pipeline-name-asdf/jobs/job-name-asdf/builds/name-asdf <<"))
		})
	})

	Describe("response", func() {
		var request Request
		var message string

		BeforeEach(func() {
			request = Request{
				Source: Source{
					Server:   "chat.freenode.net",
					Port:     7070,
					Channel:  "#random",
					User:     "randobot1337",
					Password: "secretsecret",
				},
				Params: Params{DryRun: true},
			}

			os.Setenv("BUILD_ID", "id-123")
			os.Setenv("BUILD_NAME", "name-asdf")
			os.Setenv("BUILD_JOB_NAME", "job-name-asdf")
			os.Setenv("BUILD_PIPELINE_NAME", "pipeline-name-asdf")
			os.Setenv("BUILD_TEAM_NAME", "team-name-asdf")
			os.Setenv("ATC_EXTERNAL_URL", "https://ci.example.com")

			message = "this is a message"
		})

		It("contains version", func() {
			response := BuildResponse(&request, message)
			Expect(response.Version.Ref).To(Equal("none"))
		})

		It("contains specific metadata", func() {
			response := BuildResponse(&request, message)
			Expect(response.Metadata).To(Equal([]Metadatum{
				Metadatum{"host", "chat.freenode.net:7070"},
				Metadatum{"channel", "#random"},
				Metadatum{"user", "randobot1337"},
				Metadatum{"message", "this is a message"},
			}))
		})

		It("does not contains metadata `password`", func() {
			response := BuildResponse(&request, message)
			for _, metadatum := range response.Metadata {
				Expect(metadatum.Name).To(Not(MatchRegexp(`password`)))
			}
		})
	})
})
