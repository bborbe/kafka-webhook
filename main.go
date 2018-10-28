// Copyright (c) 2018 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	flag "github.com/bborbe/flagenv"
	"github.com/bborbe/kafka-webhook/webhook"
	"github.com/golang/glog"
)

func main() {
	defer glog.Flush()
	glog.CopyStandardLogTo("info")
	runtime.GOMAXPROCS(runtime.NumCPU())

	app := &webhook.App{}
	flag.IntVar(&app.Port, "port", 9005, "port to listen")
	flag.StringVar(&app.KafkaBrokers, "kafka-brokers", "", "kafka brokers")
	flag.StringVar(&app.KafkaGroup, "kafka-group", "", "kafka consumer group")
	flag.StringVar(&app.KafkaTopic, "kafka-topic", "", "kafka topic")
	flag.StringVar(&app.HookMethod, "hook-method", http.MethodPost, "used to send data")
	flag.StringVar(&app.HookURL, "hook-url", "", "url send data to")
	flag.DurationVar(&app.RetryDelay, "retry-delay", time.Second, "amount * attempt of time to wait between retry delivery")
	flag.IntVar(&app.RetryLimit, "retry-limit", -1, "amount of retries before message is skip")

	_ = flag.Set("logtostderr", "true")
	flag.Parse()

	glog.V(0).Infof("Parameter HookMethod: %s", app.HookMethod)
	glog.V(0).Infof("Parameter HookURL: %s", app.HookURL)
	glog.V(0).Infof("Parameter KafkaBrokers: %s", app.KafkaBrokers)
	glog.V(0).Infof("Parameter KafkaGroup: %s", app.KafkaGroup)
	glog.V(0).Infof("Parameter KafkaTopic: %s", app.KafkaTopic)
	glog.V(0).Infof("Parameter Port: %d", app.Port)
	glog.V(0).Infof("Parameter RetryDelay: %v", app.RetryDelay)
	glog.V(0).Infof("Parameter RetryLimit: %d", app.RetryLimit)

	err := app.Validate()
	if err != nil {
		glog.Exit(err)
	}

	ctx := contextWithSig(context.Background())

	glog.V(0).Infof("app started")
	if err := app.Run(ctx); err != nil {
		glog.Exitf("app failed: %+v", err)
	}
	glog.V(0).Infof("app finished")
}

func contextWithSig(ctx context.Context) context.Context {
	ctxWithCancel, cancel := context.WithCancel(ctx)
	go func() {
		defer cancel()

		signalCh := make(chan os.Signal, 1)
		signal.Notify(signalCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

		select {
		case <-signalCh:
		case <-ctx.Done():
		}
	}()

	return ctxWithCancel
}
