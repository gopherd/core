package httputil

import (
	"net/url"
	"testing"
)

func TestSHA256Signer(t *testing.T) {
	const (
		key       = "abcdef123456"
		keyField  = "private_key"
		signField = "sign"
	)
	args := url.Values{
		"y": {"2"},
		"z": {"3"},
		"x": {"1"},
	}
	sign, err := SHA256Signer()(key, keyField, signField, args)
	if err != nil {
		t.Fatalf("sign error: %v", err)
		return
	}
	const expected = "21dca7ab28c84ef4566021ead8f9bb876e3d2945d2c7994c809e26ee236b531c"
	if sign != expected {
		t.Errorf("sign expected %s, but got %s", expected, sign)
		return
	}

	args.Set(signField, sign)
	if err := VerifySign(SHA256Signer(), key, keyField, signField, args); err != nil {
		t.Errorf("verify sign failure: %v", err)
	}
}
