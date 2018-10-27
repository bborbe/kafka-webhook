// Copyright (c) 2018 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Kafka Webhook", func() {
	It("Compiles", func() {
		var err error
		_, err = gexec.Build("github.com/bborbe/kafka-webhook")
		Expect(err).NotTo(HaveOccurred())
	})
})

func TestKafkaK8sVersionCollector(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Kafka Webhook Suite")
}
