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
	glog.V(3).Infof("complete %s request to %s in %d ms with status %d", req.Method, req.URL.String(), int(duration*1000), resp.StatusCode)
	httpRequestDurationSummary.Observe(duration)
	return resp, err
}
