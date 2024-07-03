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
	AuthSession *sess.AuthorisedSession
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
	client.AuthSession.Logger = logger
	return client
}

func NewTlsSkipCustomLogger(api cfg.APIServer, auth cfg.AuthServer, logger sess.SessionLog) *Client {
	client := NewTlsSkip(api, auth)
	return client
}

func NewTlsSkip(api cfg.APIServer, auth cfg.AuthServer) *Client {
	return NewClient(buildClient(), api, auth)
}

func NewClient(cl *http.Client, api cfg.APIServer, auth cfg.AuthServer) *Client {
	client := Client{}

	asession := sess.NewAuthorisedSession(cl, api, auth)
	client.AuthSession = &asession
	return &client
}

func HasAuthSession(client *Client) bool {
	return client.AuthSession != nil
}

func SetAuthToken(client *Client, token string) error {
	if client == nil {
		return errors.New("client was nil")
	}
	if HasAuthSession(client) {
		sess.SetAToken(client.AuthSession, token)
		return nil
	}
	return errors.New("auth session was nil")
}
