// Copyright (c) 2018 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webhook

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/bborbe/run"
	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/seibert-media/go-kafka/consumer"
)

type App struct {
	HookMethod   string
	HookURL      string
	KafkaBrokers string
	KafkaGroup   string
	KafkaTopic   string
	Port         int
	RetryDelay   time.Duration
	RetryLimit   int
}

func (a *App) Validate() error {
	if a.Port <= 0 {
		return errors.New("Port invalid")
	}
	if a.KafkaBrokers == "" {
		return errors.New("KafkaBrokers missing")
	}
	if a.KafkaTopic == "" {
		return errors.New("KafkaTopic missing")
	}
	if a.KafkaGroup == "" {
		return errors.New("KafkaGroup missing")
	}
	if a.HookURL == "" {
		return errors.New("Url missing")
	}
	if a.HookMethod == "" {
		return errors.New("HookMethod missing")
	}
	return nil
}

func (a *App) Run(ctx context.Context) error {

	return run.CancelOnFirstFinish(ctx, a.RunConsumer, a.RunServer)
}

func (a *App) RunServer(ctx context.Context) error {
	router := mux.NewRouter()
	router.Path("/healthz").HandlerFunc(a.HealthCheck)
	router.Path("/readiness").HandlerFunc(a.ReadinessCheck)
	router.Path("/metrics").Handler(promhttp.Handler())
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", a.Port),
		Handler: router,
	}
	go func() {
		select {
		case <-ctx.Done():
			if err := server.Shutdown(ctx); err != nil {
				glog.Warningf("shutdown failed: %v", err)
			}
		}
	}()
	return server.ListenAndServe()
}

func (a *App) RunConsumer(ctx context.Context) error {
	consumer := &consumer.OffsetConsumer{
		KafkaBrokers: a.KafkaBrokers,
		KafkaTopic:   a.KafkaTopic,
		KafkaGroup:   a.KafkaGroup,
		MessageHandler: &RetryMessageHandler{
			MaxRetry:           a.RetryLimit,
			WaitBetweenRetries: a.RetryDelay,
			MessageHandler: &PostMessageHandler{
				Timeout: 10 * time.Second,
				Url:     a.HookURL,
				Method:  a.HookMethod,
				HttpClient: &HttpClientMetrics{
					HttpClient: http.DefaultClient,
				},
			},
		},
	}
	return consumer.Consume(ctx)
}

func (a *App) HealthCheck(resp http.ResponseWriter, req *http.Request) {
	resp.WriteHeader(200)
	fmt.Fprintf(resp, "ok")
}

func (a *App) ReadinessCheck(resp http.ResponseWriter, req *http.Request) {
	resp.WriteHeader(200)
	fmt.Fprintf(resp, "ok")
}
