// Copyright (c) 2018 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webhook_test

import (
	"context"
	"net/http"

	"github.com/Shopify/sarama"
	"github.com/bborbe/kafka-webhook/mocks"
	"github.com/bborbe/kafka-webhook/webhook"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("PostMessageHandler", func() {
	var messageHandler *webhook.PostMessageHandler
	var httpClient *mocks.HttpClient
	BeforeEach(func() {
		httpClient = &mocks.HttpClient{}
		httpClient.DoReturns(&http.Response{
			StatusCode: 200,
		}, nil)
		messageHandler = &webhook.PostMessageHandler{
			Url:        "http://www.example.com",
			Method:     http.MethodPost,
			HttpClient: httpClient,
		}
	})

	It("returns no error", func() {
		err := messageHandler.ConsumeMessage(context.Background(), &sarama.ConsumerMessage{})
		Expect(err).To(BeNil())
	})
})
