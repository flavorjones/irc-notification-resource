package irc_test

import (
	"bytes"
	"encoding/json"
	"os"

	. "github.com/flavorjones/irc-notification-resource/pkg/irc"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

type GenericMap map[string]interface{}

var _ = Describe("Out", func() {
	//
	//  make it easy to generate JSON strings that are set up within each Describe block
	//
	var sourceMap, paramsMap GenericMap
	var requestJson = func() string {
		requestMap := make(GenericMap)
		requestMap["source"] = sourceMap
		requestMap["params"] = paramsMap
		message, _ := json.Marshal(requestMap)
		return string(message)
	}
	var minimalRequestJson = func() string {
		delete(sourceMap, "usetls")
		delete(sourceMap, "join")
		delete(sourceMap, "debug")
		delete(paramsMap, "dry_run")
		return requestJson()
	}
	var messageFileRequestJson = func() string {
		delete(paramsMap, "message")
		paramsMap["message_file"] = "test/one-line-message-file.txt"
		return minimalRequestJson()
	}

	BeforeEach(func() {
		sourceMap = make(GenericMap)
		paramsMap = make(GenericMap)
		sourceMap["server"] = "chat.freenode.net"
		sourceMap["port"] = 7070
		sourceMap["channel"] = "#random"
		sourceMap["user"] = "randobot1337"
		sourceMap["password"] = "secretsecret"
		sourceMap["usetls"] = true
		sourceMap["join"] = false
		sourceMap["debug"] = false
		paramsMap["message"] = "foo $BUILD_ID"
		paramsMap["dry_run"] = false
	})

	Describe("VERSION", func() {
		It("provides a version string", func() {
			Expect(VERSION).To(MatchRegexp(`^v\d+\.\d+\.\d+$`))
		})
	})

	Describe("test setup", func() {
		Describe("test request json", func() {
			It("is as expected", func() {
				Expect(requestJson()).To(Equal(`{"params":{"dry_run":false,"message":"foo $BUILD_ID"},"source":{"channel":"#random","debug":false,"join":false,"password":"secretsecret","port":7070,"server":"chat.freenode.net","user":"randobot1337","usetls":true}}`))
			})

			Context("without a key", func() {
				BeforeEach(func() {
					delete(sourceMap, "channel")
				})
				It("is as expected", func() {
					Expect(requestJson()).To(Equal(`{"params":{"dry_run":false,"message":"foo $BUILD_ID"},"source":{"debug":false,"join":false,"password":"secretsecret","port":7070,"server":"chat.freenode.net","user":"randobot1337","usetls":true}}`))
				})
			})
		})

		Describe("test minimal request json", func() {
			It("is as expected", func() {
				Expect(minimalRequestJson()).To(Equal(`{"params":{"message":"foo $BUILD_ID"},"source":{"channel":"#random","password":"secretsecret","port":7070,"server":"chat.freenode.net","user":"randobot1337"}}`))
			})

			Context("without a key", func() {
				BeforeEach(func() {
					delete(sourceMap, "channel")
				})
				It("is as expected", func() {
					Expect(minimalRequestJson()).To(Equal(`{"params":{"message":"foo $BUILD_ID"},"source":{"password":"secretsecret","port":7070,"server":"chat.freenode.net","user":"randobot1337"}}`))
				})
			})
		})

		Describe("test message_file request json", func() {
			It("is as expected", func() {
				Expect(messageFileRequestJson()).To(Equal(`{"params":{"message_file":"test/one-line-message-file.txt"},"source":{"channel":"#random","password":"secretsecret","port":7070,"server":"chat.freenode.net","user":"randobot1337"}}`))
			})
		})
	})

	Describe("ParseAndCheckRequest()", func() {
		It("returns correct Source values", func() {
			request, error := ParseAndCheckRequest(bytes.NewBufferString(requestJson()))
			Expect(error).To(BeNil())
			Expect(request.Source).To(MatchAllFields(Fields{
				"Server":   Equal("chat.freenode.net"),
				"Port":     Equal(7070),
				"Channel":  Equal("#random"),
				"User":     Equal("randobot1337"),
				"Password": Equal("secretsecret"),
				"UseTLS":   Equal(true),
				"Join":     Equal(false),
				"Debug":    Equal(false),
			}))
		})

		Describe("required Source property", func() {
			Describe("`server`", func() {
				Context("when not present", func() {
					BeforeEach(func() {
						delete(sourceMap, "server")
					})
					It("errors", func() {
						_, error := ParseAndCheckRequest(bytes.NewBufferString(minimalRequestJson()))
						Expect(error.Error()).To(MatchRegexp(`No server was provided`))
					})
				})
			})

			Describe("`port`", func() {
				Context("when not present", func() {
					BeforeEach(func() {
						delete(sourceMap, "port")
					})
					It("errors", func() {
						_, error := ParseAndCheckRequest(bytes.NewBufferString(minimalRequestJson()))
						Expect(error.Error()).To(MatchRegexp(`No port was provided`))
					})
				})
			})

			Describe("`channel`", func() {
				Context("when not present", func() {
					BeforeEach(func() {
						delete(sourceMap, "channel")
					})
					It("errors", func() {
						_, error := ParseAndCheckRequest(bytes.NewBufferString(minimalRequestJson()))
						Expect(error.Error()).To(MatchRegexp(`No channel was provided`))
					})
				})
			})

			Describe("`user`", func() {
				Context("when not present", func() {
					BeforeEach(func() {
						delete(sourceMap, "user")
					})
					It("errors", func() {
						_, error := ParseAndCheckRequest(bytes.NewBufferString(minimalRequestJson()))
						Expect(error.Error()).To(MatchRegexp(`No user was provided`))
					})
				})
			})

			Describe("`password`", func() {
				Context("when not present", func() {
					BeforeEach(func() {
						delete(sourceMap, "password")
					})
					It("errors", func() {
						_, error := ParseAndCheckRequest(bytes.NewBufferString(minimalRequestJson()))
						Expect(error.Error()).To(MatchRegexp(`No password was provided`))
					})
				})
			})
		})

		Describe("optional Source property", func() {
			Describe("`usetls`", func() {
				Context("when not present", func() {
					It("defaults to true", func() {
						request, error := ParseAndCheckRequest(bytes.NewBufferString(minimalRequestJson()))
						Expect(error).To(BeNil())
						Expect(request.Source).To(MatchFields(IgnoreExtras, Fields{"UseTLS": BeTrue()}))
					})
				})

				Context("when set to true", func() {
					BeforeEach(func() {
						sourceMap["usetls"] = true
					})
					It("is true", func() {
						request, error := ParseAndCheckRequest(bytes.NewBufferString(requestJson()))
						Expect(error).To(BeNil())
						Expect(request.Source).To(MatchFields(IgnoreExtras, Fields{"UseTLS": BeTrue()}))
					})
				})

				Context("when set to false", func() {
					BeforeEach(func() {
						sourceMap["usetls"] = false
					})
					It("is false", func() {
						request, error := ParseAndCheckRequest(bytes.NewBufferString(requestJson()))
						Expect(error).To(BeNil())
						Expect(request.Source).To(MatchFields(IgnoreExtras, Fields{"UseTLS": BeFalse()}))
					})
				})
			})

			Describe("`join`", func() {
				Context("when not present", func() {
					It("defaults to false", func() {
						request, error := ParseAndCheckRequest(bytes.NewBufferString(minimalRequestJson()))
						Expect(error).To(BeNil())
						Expect(request.Source).To(MatchFields(IgnoreExtras, Fields{"Join": BeFalse()}))
					})
				})

				Context("when set to true", func() {
					BeforeEach(func() {
						sourceMap["join"] = true
					})
					It("is true", func() {
						request, error := ParseAndCheckRequest(bytes.NewBufferString(requestJson()))
						Expect(error).To(BeNil())
						Expect(request.Source).To(MatchFields(IgnoreExtras, Fields{"Join": BeTrue()}))
					})
				})

				Context("when set to false", func() {
					BeforeEach(func() {
						sourceMap["join"] = false
					})
					It("is false", func() {
						request, error := ParseAndCheckRequest(bytes.NewBufferString(requestJson()))
						Expect(error).To(BeNil())
						Expect(request.Source).To(MatchFields(IgnoreExtras, Fields{"Join": BeFalse()}))
					})
				})
			})

			Describe("`debug`", func() {
				Context("when not present", func() {
					It("defaults to false", func() {
						request, error := ParseAndCheckRequest(bytes.NewBufferString(minimalRequestJson()))
						Expect(error).To(BeNil())
						Expect(request.Source).To(MatchFields(IgnoreExtras, Fields{"Debug": BeFalse()}))
					})
				})

				Context("when set to true", func() {
					BeforeEach(func() {
						sourceMap["debug"] = true
					})
					It("is true", func() {
						request, error := ParseAndCheckRequest(bytes.NewBufferString(requestJson()))
						Expect(error).To(BeNil())
						Expect(request.Source).To(MatchFields(IgnoreExtras, Fields{"Debug": BeTrue()}))
					})
				})

				Context("when set to false", func() {
					BeforeEach(func() {
						sourceMap["debug"] = false
					})
					It("is false", func() {
						request, error := ParseAndCheckRequest(bytes.NewBufferString(requestJson()))
						Expect(error).To(BeNil())
						Expect(request.Source).To(MatchFields(IgnoreExtras, Fields{"Debug": BeFalse()}))
					})
				})
			})
		})

		Context("when given a message", func() {
			It("returns correct Params values", func() {
				request, error := ParseAndCheckRequest(bytes.NewBufferString(minimalRequestJson()))
				Expect(error).To(BeNil())
				Expect(request.Params.MessageFile).To(Equal(""))
				Expect(request.Params).To(MatchFields(IgnoreExtras, Fields{
					"Message": Equal("foo $BUILD_ID"),
				}))
			})
		})

		Context("when given a message_file", func() {
			It("returns correct Params values", func() {
				request, error := ParseAndCheckRequest(bytes.NewBufferString(messageFileRequestJson()))
				Expect(error).To(BeNil())
				Expect(request.Params.Message).To(Equal(""))
				Expect(request.Params).To(MatchFields(IgnoreExtras, Fields{
					"MessageFile": Equal("test/one-line-message-file.txt"),
				}))
			})
		})

		Describe("required Params property", func() {
			Describe("`message or message_file`", func() {
				Context("when neither is present", func() {
					BeforeEach(func() {
						delete(paramsMap, "message")
						delete(paramsMap, "message_file")
					})
					It("errors", func() {
						_, error := ParseAndCheckRequest(bytes.NewBufferString(minimalRequestJson()))
						Expect(error.Error()).To(MatchRegexp(`Either message or message_file must be set`))
					})
				})
			})
		})

		Describe("optional Params property", func() {
			Describe("`dry_run`", func() {
				Context("when not set", func() {
					It("defaults to false", func() {
						request, error := ParseAndCheckRequest(bytes.NewBufferString(minimalRequestJson()))
						Expect(error).To(BeNil())
						Expect(request.Params).To(MatchFields(IgnoreExtras, Fields{"DryRun": BeFalse()}))
					})
				})

				Context("when set to true", func() {
					BeforeEach(func() {
						paramsMap["dry_run"] = true
					})
					It("is true", func() {
						request, error := ParseAndCheckRequest(bytes.NewBufferString(requestJson()))
						Expect(error).To(BeNil())
						Expect(request.Params).To(MatchFields(IgnoreExtras, Fields{"DryRun": BeTrue()}))
					})
				})

				Context("when set to false", func() {
					BeforeEach(func() {
						paramsMap["dry_run"] = false
					})
					It("is false", func() {
						request, error := ParseAndCheckRequest(bytes.NewBufferString(requestJson()))
						Expect(error).To(BeNil())
						Expect(request.Params).To(MatchFields(IgnoreExtras, Fields{"DryRun": BeFalse()}))
					})
				})
			})
		})
	})

	Describe("ExpandMessage()", func() {
		BeforeEach(func() {
			os.Setenv("BUILD_ID", "id-123")
			os.Setenv("BUILD_NAME", "name-asdf")
			os.Setenv("BUILD_JOB_NAME", "job-name-asdf")
			os.Setenv("BUILD_PIPELINE_NAME", "pipeline-name-asdf")
			os.Setenv("BUILD_TEAM_NAME", "team-name-asdf")
			os.Setenv("ATC_EXTERNAL_URL", "https://ci.example.com")
		})

		It("expands environment variables", func() {
			message := ">> $BUILD_ID <<"
			expanded := ExpandMessage(message)
			Expect(expanded).To(Equal(">> id-123 <<"))
		})

		It("expands BUILD_URL pseudo-metadata", func() {
			message := ">> $BUILD_URL <<"
			expanded := ExpandMessage(message)
			Expect(expanded).To(Equal(">> https://ci.example.com/teams/team-name-asdf/pipelines/pipeline-name-asdf/jobs/job-name-asdf/builds/name-asdf <<"))
		})
	})

	Describe("BuildResponse()", func() {
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
					UseTLS:   true,
					Join:     false,
					Debug:    false,
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

		Describe("returned Response", func() {
			It("contains version", func() {
				response := BuildResponse(&request, message)
				Expect(response.Version.Ref).To(Equal("none"))
			})

			It("contains specific metadata", func() {
				response := BuildResponse(&request, message)
				Expect(response.Metadata).To(Equal([]Metadatum{
					Metadatum{"resource_version", VERSION},
					Metadatum{"host", "chat.freenode.net:7070"},
					Metadatum{"channel", "#random"},
					Metadatum{"user", "randobot1337"},
					Metadatum{"usetls", "true"},
					Metadatum{"join", "false"},
					Metadatum{"debug", "false"},
					Metadatum{"message", "this is a message"},
					Metadatum{"dry_run", "true"},
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
})
