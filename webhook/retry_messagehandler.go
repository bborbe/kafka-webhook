// Copyright (c) 2018 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webhook

import (
	"context"
	"time"

	"github.com/Shopify/sarama"
	"github.com/golang/glog"
	"github.com/pkg/errors"
)

// RetryMessageHandler calls the given MessageHandler until no errors occurred, MaxRetry reached or Context is done.
type RetryMessageHandler struct {
	// MessageHandler to call
	MessageHandler MessageHandler
	// MaxRetry is the amount of retries. Negative number let retry forever
	MaxRetry int
	// RetryDelay is the amount of time * retry counter wait before the next try
	WaitBetweenRetries time.Duration
}

// ConsumeMessage send the message to the given MessageHandler and retries if needed.
func (r *RetryMessageHandler) ConsumeMessage(ctx context.Context, msg *sarama.ConsumerMessage) error {
	glog.V(3).Infof("consume message of topic %s, partition %d and offset %d", msg.Topic, msg.Partition, msg.Offset)
	counter := 0
	for {
		counter++
		err := r.MessageHandler.ConsumeMessage(ctx, msg)
		if err == nil {
			glog.V(3).Infof("consume message successful")
			return nil
		}
		glog.V(3).Infof("message handler returned error => retry")
		if r.MaxRetry >= 0 && counter > r.MaxRetry {
			return errors.Wrapf(err, "max retries reached")
		}
		wait := r.WaitBetweenRetries * time.Duration(counter)
		glog.V(1).Infof("handle message failed %d times => retry in %v", counter, wait)
		if wait > 0 {
			select {
			case <-ctx.Done():
				return nil
			case <-time.NewTimer(wait).C:
				glog.V(3).Infof("wait for %v completed", wait)
			}
		}
	}
}
