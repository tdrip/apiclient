package util

import (
	"crypto/tls"
	"net/http"
	"time"
)

func BuildNoTLSClient() *http.Client {
	transport := &http.Transport{Proxy: http.ProxyFromEnvironment, TLSClientConfig: &tls.Config{Renegotiation: tls.RenegotiateOnceAsClient, InsecureSkipVerify: true}, TLSNextProto: map[string]func(authority string, c *tls.Conn) http.RoundTripper{}}
	return &http.Client{
		Timeout:   time.Second * 10,
		Transport: transport,
	}
}

func RequestIsSuccessful(code int) bool {
	return (code >= 200 && code < 300)
}
