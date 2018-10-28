// Copyright (c) 2018 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webhook

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/Shopify/sarama"
	"github.com/pkg/errors"
)

const (
	KeyField       = "X-Message-Key"
	TopicField     = "X-Message-Topic"
	OffsetField    = "X-Message-Offset"
	PartitionField = "X-Message-Partition"
	SignaturField  = "X-Signature"
)

type RequestCoding struct {
	Url    string
	Method string
	Signer interface {
		Compare(content []byte, sign string) (bool, error)
		Sign(content []byte) string
	}
}

func (r *RequestCoding) Encode(msg *sarama.ConsumerMessage) (*http.Request, error) {
	req, err := http.NewRequest(r.Method, r.Url, bytes.NewBuffer(msg.Value))
	if err != nil {
		return nil, errors.Wrap(err, "build request failed")
	}
	req.Header.Add(KeyField, base64.StdEncoding.EncodeToString(msg.Key))
	req.Header.Add(TopicField, msg.Topic)
	req.Header.Add(OffsetField, strconv.FormatInt(msg.Offset, 10))
	req.Header.Add(PartitionField, strconv.FormatInt(int64(msg.Partition), 10))
	req.Header.Add(SignaturField, r.Signer.Sign(msg.Value))
	return req, nil
}

func (r *RequestCoding) Decode(req *http.Request) (*sarama.ConsumerMessage, error) {
	content, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, errors.Wrap(err, "read body failed")
	}
	defer req.Body.Close()
	equal, err := r.Signer.Compare(content, req.Header.Get(SignaturField))
	if err != nil {
		return nil, fmt.Errorf("compare signs failed")
	}
	if !equal {
		return nil, fmt.Errorf("check secret failed")
	}
	key, err := base64.StdEncoding.DecodeString(req.Header.Get(KeyField))
	if err != nil {
		return nil, errors.Wrap(err, "decode key failed")
	}
	offset, err := strconv.Atoi(req.Header.Get(OffsetField))
	if err != nil {
		return nil, errors.Wrap(err, "decode offset failed")
	}
	partition, err := strconv.Atoi(req.Header.Get(PartitionField))
	if err != nil {
		return nil, errors.Wrap(err, "decode partition failed")
	}
	return &sarama.ConsumerMessage{
		Value:     content,
		Key:       key,
		Topic:     req.Header.Get(TopicField),
		Offset:    int64(offset),
		Partition: int32(partition),
	}, nil
}
