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

type MessageHandler struct {
	Url        string
	Method     string
	HttpClient HttpClient
}

func (m *MessageHandler) ConsumeMessage(ctx context.Context, msg *sarama.ConsumerMessage) error {
	req, err := http.NewRequest(m.Method, m.Url, bytes.NewBuffer(msg.Value))
	if err != nil {
		return errors.Wrap(err, "build request failed")
	}
	req.Header.Add("X-Message-Key", base64.StdEncoding.EncodeToString(msg.Key))
	req.Header.Add("X-Message-Topic", msg.Topic)
	req.Header.Add("X-Message-Offset", strconv.FormatInt(msg.Offset, 10))
	req.Header.Add("X-Message-Partition", strconv.FormatInt(int64(msg.Partition), 10))
	req.Header.Add("X-Message-Timestamp", msg.Timestamp.Format(time.RFC3339))
	req.Header.Add("X-Message-Block-Timestamp", msg.BlockTimestamp.Format(time.RFC3339))

	resp, err := m.HttpClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "perform request failed")
	}
	if resp.StatusCode/100 != 2 {
		return errors.Wrap(err, "status != 2xx")
	}
	return nil
}
