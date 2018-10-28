// Copyright (c) 2018 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webhook_test

import (
	"context"
	"time"

	"github.com/Shopify/sarama"
	"github.com/bborbe/kafka-webhook/mocks"
	"github.com/bborbe/kafka-webhook/webhook"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
)

var _ = Describe("RetryMessageHandler", func() {
	var retryMessageHandler *webhook.RetryMessageHandler
	var messageHandler *mocks.MessageHandler
	BeforeEach(func() {
		messageHandler = &mocks.MessageHandler{}
		retryMessageHandler = &webhook.RetryMessageHandler{
			MessageHandler:     messageHandler,
			MaxRetry:           -1,
			WaitBetweenRetries: time.Nanosecond,
		}
	})

	It("calls messagehandler", func() {
		ctx := context.Background()
		message := &sarama.ConsumerMessage{}

		err := retryMessageHandler.ConsumeMessage(ctx, message)
		Expect(err).To(BeNil())
		Expect(messageHandler.ConsumeMessageCallCount()).To(Equal(1))
		argCtx, argMessage := messageHandler.ConsumeMessageArgsForCall(0)
		Expect(argCtx).To(Equal(ctx))
		Expect(argMessage).To(Equal(message))
	})

	It("calls messagehandler until no error", func() {
		messageHandler.ConsumeMessageReturnsOnCall(0, errors.New("banana"))
		messageHandler.ConsumeMessageReturnsOnCall(1, errors.New("banana"))
		messageHandler.ConsumeMessageReturnsOnCall(2, errors.New("banana"))
		messageHandler.ConsumeMessageReturnsOnCall(3, errors.New("banana"))
		messageHandler.ConsumeMessageReturnsOnCall(4, nil)

		err := retryMessageHandler.ConsumeMessage(context.Background(), &sarama.ConsumerMessage{})
		Expect(err).To(BeNil())
		Expect(messageHandler.ConsumeMessageCallCount()).To(Equal(5))
	})

	It("calls messagehandler until max retries reached no error", func() {
		retryMessageHandler.MaxRetry = 9
		messageHandler.ConsumeMessageReturns(errors.New("banana"))

		err := retryMessageHandler.ConsumeMessage(context.Background(), &sarama.ConsumerMessage{})
		Expect(err).NotTo(BeNil())
		Expect(messageHandler.ConsumeMessageCallCount()).To(Equal(retryMessageHandler.MaxRetry + 1))
	})
})
