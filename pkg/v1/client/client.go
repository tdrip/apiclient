package client

import (
	"errors"
	"net/http"

	cfg "github.com/tdrip/apiclient/pkg/v1/cfg"
	sess "github.com/tdrip/apiclient/pkg/v1/session"
	utils "github.com/tdrip/apiclient/pkg/v1/utils"
)

type Client struct {
	Session sess.Session
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
	return NewClient(utils.BuildNoTLSClient(), api, auth)
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
