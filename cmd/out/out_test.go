package main_test

import (
	"bytes"

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
})
