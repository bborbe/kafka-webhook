// Copyright (c) 2018 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webhook_test

import (
	"encoding/base64"
	"net/http"

	"github.com/Shopify/sarama"
	"github.com/bborbe/kafka-webhook/webhook"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/format"
)

var _ = Describe("RequestCoding", func() {
	format.TruncatedDiff = false
	var requestCoding *webhook.RequestCoding
	var msg *sarama.ConsumerMessage
	var offset int64 = 123
	var partition int32 = 23
	key := []byte("key")
	value := []byte("value")
	topic := "myTopic"
	BeforeEach(func() {
		requestCoding = &webhook.RequestCoding{
			Url:    "http://example.com/hook",
			Method: http.MethodPost,
			Signer: &webhook.Signer{
				Secret: "secret",
			},
		}
		msg = &sarama.ConsumerMessage{
			Offset:    offset,
			Partition: partition,
			Key:       key,
			Value:     value,
			Topic:     topic,
		}
	})
	It("set topic as header", func() {
		req, err := requestCoding.Encode(msg)
		Expect(err).To(BeNil())
		Expect(req.Header.Get(webhook.TopicField)).To(Equal(topic))
	})
	It("set key as header", func() {
		req, err := requestCoding.Encode(msg)
		Expect(err).To(BeNil())
		value, err := base64.StdEncoding.DecodeString(req.Header.Get(webhook.KeyField))
		Expect(err).To(BeNil())
		Expect(value).To(Equal(key))
	})
	It("set offset as header", func() {
		req, err := requestCoding.Encode(msg)
		Expect(err).To(BeNil())
		Expect(req.Header.Get(webhook.OffsetField)).To(Equal("123"))
	})
	It("set partition as header", func() {
		req, err := requestCoding.Encode(msg)
		Expect(err).To(BeNil())
		Expect(req.Header.Get(webhook.PartitionField)).To(Equal("23"))
	})
	It("set signatur as header", func() {
		req, err := requestCoding.Encode(msg)
		Expect(err).To(BeNil())
		value, ok := req.Header[webhook.SignaturField]
		Expect(ok).To(BeTrue())
		Expect(value).To(HaveLen(1))
		Expect(value[0]).To(Equal("50e03ebe65be98bb8bf11ba2c892d54c079aca2b0d3b0162769c6d757a25434f"))
	})
	It("encodes sarama message to request and back", func() {
		req, err := requestCoding.Encode(msg)
		Expect(err).To(BeNil())
		Expect(req).NotTo(BeNil())
		message, err := requestCoding.Decode(req)
		Expect(err).To(BeNil())
		Expect(message).NotTo(BeNil())
		Expect(message.Key).To(Equal(key))
		Expect(message.Value).To(Equal(value))
		Expect(message.Offset).To(Equal(offset))
		Expect(message.Partition).To(Equal(partition))
	})
})
