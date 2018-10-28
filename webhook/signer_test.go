// Copyright (c) 2018 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webhook_test

import (
	"github.com/bborbe/kafka-webhook/webhook"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Signer", func() {
	It("", func() {
		signer := &webhook.Signer{
			Secret: "s3cr3t",
		}
		content := []byte("banana")
		sign := signer.Sign(content)
		equal, err := signer.Compare(content, sign)
		Expect(err).To(BeNil())
		Expect(equal).To(BeTrue())
	})
})
