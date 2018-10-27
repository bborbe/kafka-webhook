// Copyright (c) 2018 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webhook

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bborbe/run"
	"github.com/golang/glog"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/seibert-media/go-kafka/consumer"
)

type App struct {
	KafkaBrokers string
	Port         int
	KafkaTopic   string
	KafkaGroup   string
	Url          string
	Method       string
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
	if a.Url == "" {
		return errors.New("Url missing")
	}
	if a.Method == "" {
		return errors.New("Method missing")
	}
	return nil
}

func (a *App) Run(ctx context.Context) error {

	return run.CancelOnFirstFinish(ctx, a.RunConsumer, a.RunServer)
}

func (a *App) RunServer(ctx context.Context) error {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", a.Port),
		Handler: promhttp.Handler(),
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
		MessageHandler: &MessageHandler{
			Url:    a.Url,
			Method: a.Method,
			HttpClient: &HttpClientMetrics{
				HttpClient: http.DefaultClient,
			},
		},
	}
	return consumer.Consume(ctx)
}
