package client

import (
	"crypto/tls"
	"errors"
	"net/http"
	"time"

	cfg "github.com/tdrip/apiclient/pkg/v1/cfg"
	sess "github.com/tdrip/apiclient/pkg/v1/session"
)

type Client struct {
	Session sess.Session
}

func buildClient() *http.Client {
	transport := &http.Transport{Proxy: http.ProxyFromEnvironment, TLSClientConfig: &tls.Config{Renegotiation: tls.RenegotiateOnceAsClient, InsecureSkipVerify: true}, TLSNextProto: map[string]func(authority string, c *tls.Conn) http.RoundTripper{}}
	return &http.Client{
		Timeout:   time.Second * 10,
		Transport: transport,
	}
}

func NewClientCustomLogger(cl *http.Client, api cfg.APIServer, auth cfg.AuthServer, logger sess.SessionLog) *Client {
	client := NewClient(cl, api, auth)
	client.Session.Logger = logger
	return client
}

func NewTlsSkipCustomLogger(api cfg.APIServer, auth cfg.AuthServer, logger sess.SessionLog) *Client {
	client := NewTlsSkip(api, auth)
	client.Session.Logger = logger
	return client
}

func NewTlsSkip(api cfg.APIServer, auth cfg.AuthServer) *Client {
	return NewClient(buildClient(), api, auth)
}

func NewClient(cl *http.Client, api cfg.APIServer, auth cfg.AuthServer) *Client {
	client := Client{Session: sess.NewSession(cl, api, auth)}
	return &client
}

func SetAuthToken(client *Client, token string) error {
	if client == nil {
		return errors.New("client was nil")
	}
	client.Session = client.Session.UpdateAToken(token)
	return nil
}
