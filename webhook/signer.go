// Copyright (c) 2018 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webhook

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"

	"github.com/pkg/errors"
)

type Signer struct {
	Secret string
}

func (r *Signer) Compare(content []byte, sign string) (bool, error) {
	mac1 := r.signBody(content)
	mac2, err := hex.DecodeString(sign)
	if err != nil {
		return false, errors.Wrap(err, "decode sign failed")
	}
	return hmac.Equal(mac1, mac2), nil
}

func (r *Signer) Sign(content []byte) string {
	return hex.EncodeToString(r.signBody(content))
}

func (r *Signer) signBody(content []byte) []byte {
	computed := hmac.New(sha256.New, []byte(r.Secret))
	computed.Write(content)
	return []byte(computed.Sum(nil))
}
