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

func NewTlsSkip(api cfg.APIServer, auth cfg.AuthServer) (*Client, error) {
	return New(buildClient(), api, auth)
}

func New(cl *http.Client, api cfg.APIServer, auth cfg.AuthServer) (*Client, error) {
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

func SetAuthToken(client *Client, token string) error {

	if client.AuthSession != nil {
		auth := client.AuthSession
		updated := (*auth).SetAuthToken(token)
		client.AuthSession = &updated
	}

	return errors.New("Auth Session was nil")
}
