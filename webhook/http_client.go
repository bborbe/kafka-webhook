// Copyright (c) 2018 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webhook

import (
	"net/http"
	"time"

	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"
)

const namespace = "webhook"

var (
	httpRequestDurationSummary = prometheus.NewSummary(prometheus.SummaryOpts{
		Namespace: namespace,
		Name:      "http_duration",
		Help:      "http request duration",
	})
)

func init() {
	prometheus.MustRegister(
		httpRequestDurationSummary,
	)
}

//go:generate counterfeiter -o ../mocks/http_client.go --fake-name HttpClient . HttpClient
type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type HttpClientMetrics struct {
	HttpClient HttpClient
}

func (h *HttpClientMetrics) Do(req *http.Request) (*http.Response, error) {
	start := time.Now()
	resp, err := h.HttpClient.Do(req)
	duration := time.Since(start).Seconds()
	httpRequestDurationSummary.Observe(duration)
	if err != nil {
		glog.V(3).Infof("failed %s request to %s in %d ms: %v", req.Method, req.URL.String(), int(duration*1000), err)
		return nil, err
	}
	glog.V(3).Infof("complete %s request to %s in %d ms with status %d", req.Method, req.URL.String(), int(duration*1000), resp.StatusCode)
	return resp, nil
}
