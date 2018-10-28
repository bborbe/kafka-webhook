// Copyright (c) 2018 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webhook

import (
	"context"
	"net/http"
	"time"

	"github.com/Shopify/sarama"
	"github.com/pkg/errors"
)

type PostMessageHandler struct {
	HttpClient     HttpClient
	Timeout        time.Duration
	RequestBuilder interface {
		Encode(msg *sarama.ConsumerMessage) (*http.Request, error)
	}
}

func (p *PostMessageHandler) ConsumeMessage(ctx context.Context, msg *sarama.ConsumerMessage) error {
	req, err := p.RequestBuilder.Encode(msg)
	if err != nil {
		return errors.Wrap(err, "build request failed")
	}

	ctx, cancelFunc := context.WithTimeout(ctx, p.Timeout)
	defer cancelFunc()

	resp, err := p.HttpClient.Do(req.WithContext(ctx))
	if err != nil {
		return errors.Wrap(err, "perform request failed")
	}
	if resp.StatusCode/100 != 2 {
		return errors.Wrap(err, "status != 2xx")
	}
	return nil
}
