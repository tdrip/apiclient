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
	Session     *sess.Session
}

func buildClient() *http.Client {
	transport := &http.Transport{Proxy: http.ProxyFromEnvironment, TLSClientConfig: &tls.Config{Renegotiation: tls.RenegotiateOnceAsClient, InsecureSkipVerify: true}, TLSNextProto: map[string]func(authority string, c *tls.Conn) http.RoundTripper{}}
	return &http.Client{
		Timeout:   time.Second * 10,
		Transport: transport,
	}
}

func NewClientCustomLogger(cl *http.Client, api cfg.APIServer, auth cfg.AuthServer, logger sess.SessionLog) (*Client, error) {
	client, err := NewClient(cl, api, auth)
	if err == nil {
		client.Session.Logger = logger
	}
	return client, err
}

func NewTlsSkipCustomLogger(api cfg.APIServer, auth cfg.AuthServer, logger sess.SessionLog) (*Client, error) {
	client, err := NewTlsSkip(api, auth)
	if err == nil {
		client.Session.Logger = logger
	}
	return client, err
}

func NewTlsSkip(api cfg.APIServer, auth cfg.AuthServer) (*Client, error) {
	return NewClient(buildClient(), api, auth)
}

func NewClient(cl *http.Client, api cfg.APIServer, auth cfg.AuthServer) (*Client, error) {
	client := Client{}

	asession, err := sess.NewAuthorisedSession(cl, api, auth)
	if err != nil {
		return nil, err
	}
	session, err := sess.NewSession(cl, api)
	if err != nil {
		return nil, err
	}
	client.Session = &session
	client.AuthSession = &asession
	return &client, nil
}

func HasAuthSession(client *Client) bool {
	return client.AuthSession != nil
}

func SetAuthToken(client *Client, token string) error {

	if HasAuthSession(client) {
		client.AuthSession.SetAToken(token)
	}

	return errors.New("auth Session was nil")
}
