// Copyright (c) 2018 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webhook

import (
	"context"

	"github.com/Shopify/sarama"
)

//go:generate counterfeiter -o ../mocks/messagehandler.go --fake-name MessageHandler . MessageHandler
type MessageHandler interface {
	ConsumeMessage(ctx context.Context, msg *sarama.ConsumerMessage) error
}
