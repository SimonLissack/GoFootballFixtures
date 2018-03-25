package client

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("footballDataOrgClient", func() {
	var fdoClient footballDataOrgClient

	Describe("when making a request", func() {
		var canRequest error

		BeforeEach(func() {
			fdoClient = footballDataOrgClient{}
			fdoClient.lastRequest = time.Now()
		})

		Context("when not throttled", func() {
			JustBeforeEach(func() {
				fdoClient.requestCounterReset = -10
				fdoClient.requestsAvailable = 1
				canRequest = fdoClient.CanMakeRequest()
			})

			It("should be able to make a request", func() {
				Expect(canRequest).To(BeNil())
			})
		})

		Context("when throttled", func() {
			JustBeforeEach(func() {
				fdoClient.requestCounterReset = 10
				fdoClient.requestsAvailable = 0
				canRequest = fdoClient.CanMakeRequest()
			})

			It("should be able to make a request", func() {
				Expect(canRequest).ToNot(BeNil())
			})
		})

		Context("when requests are available", func() {
			JustBeforeEach(func() {
				fdoClient.requestCounterReset = 10
				fdoClient.requestsAvailable = 1
				canRequest = fdoClient.CanMakeRequest()
			})

			It("should be able to make a request", func() {
				Expect(canRequest).To(BeNil())
			})
		})
	})
})
