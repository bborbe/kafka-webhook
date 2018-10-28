// Copyright (c) 2018 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webhook

import (
	"bytes"
	"context"
	"encoding/base64"
	"net/http"
	"strconv"
	"time"

	"github.com/Shopify/sarama"
	"github.com/pkg/errors"
)

type PostMessageHandler struct {
	Url        string
	Method     string
	HttpClient HttpClient
	Timeout    time.Duration
}

func (m *PostMessageHandler) ConsumeMessage(ctx context.Context, msg *sarama.ConsumerMessage) error {
	req, err := http.NewRequest(m.Method, m.Url, bytes.NewBuffer(msg.Value))
	if err != nil {
		return errors.Wrap(err, "build request failed")
	}
	req.Header.Add("X-Message-Key", base64.StdEncoding.EncodeToString(msg.Key))
	req.Header.Add("X-Message-Topic", msg.Topic)
	req.Header.Add("X-Message-Offset", strconv.FormatInt(msg.Offset, 10))
	req.Header.Add("X-Message-Partition", strconv.FormatInt(int64(msg.Partition), 10))

	ctx, cancelFunc := context.WithTimeout(ctx, m.Timeout)
	defer cancelFunc()

	resp, err := m.HttpClient.Do(req.WithContext(ctx))
	if err != nil {
		return errors.Wrap(err, "perform request failed")
	}
	if resp.StatusCode/100 != 2 {
		return errors.Wrap(err, "status != 2xx")
	}
	return nil
}
