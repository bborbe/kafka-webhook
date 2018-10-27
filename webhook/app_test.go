// Copyright (c) 2018 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webhook_test

import (
	"github.com/bborbe/kafka-webhook/webhook"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Version App", func() {
	var app *webhook.App
	BeforeEach(func() {
		app = &webhook.App{
			Port:         1337,
			KafkaBrokers: "kafka:9092",
			KafkaTopic:   "my-topic",
			KafkaGroup:   "my-group",
			Url:          "http://www.example.com",
		}
	})
	It("Validate without error", func() {
		Expect(app.Validate()).NotTo(HaveOccurred())
	})
	It("Validate returns error if port is 0", func() {
		app.Port = 0
		Expect(app.Validate()).To(HaveOccurred())
	})
	It("Validate returns error KafkaBrokers is empty", func() {
		app.KafkaBrokers = ""
		Expect(app.Validate()).To(HaveOccurred())
	})
	It("Validate returns error KafkaTopic is empty", func() {
		app.KafkaTopic = ""
		Expect(app.Validate()).To(HaveOccurred())
	})
	It("Validate returns error KafkaGroup is empty", func() {
		app.KafkaGroup = ""
		Expect(app.Validate()).To(HaveOccurred())
	})
	It("Validate returns error Url is empty", func() {
		app.Url = ""
		Expect(app.Validate()).To(HaveOccurred())
	})
})
