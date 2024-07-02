package client

import (
	"crypto/tls"
	"net/http"
	"time"

	api "github.com/tdrip/apiclient/pkg/v1/api"
	cfg "github.com/tdrip/apiclient/pkg/v1/cfg"
	sess "github.com/tdrip/apiclient/pkg/v1/session"
)

type Client struct {
	Session *sess.AuthorisedSession
	APIs    map[string]api.API
}

func buildClient() *http.Client {
	transport := &http.Transport{Proxy: http.ProxyFromEnvironment, TLSClientConfig: &tls.Config{Renegotiation: tls.RenegotiateOnceAsClient, InsecureSkipVerify: true}, TLSNextProto: map[string]func(authority string, c *tls.Conn) http.RoundTripper{}}
	return &http.Client{
		Timeout:   time.Second * 10,
		Transport: transport,
	}
}

func NewTlsSkip(authtoken string, api cfg.APIServer, auth cfg.AuthServer) (*Client, error) {
	return New(buildClient(), authtoken, api, auth)
}

func New(cl *http.Client, authtoken string, api cfg.APIServer, auth cfg.AuthServer) (*Client, error) {
	client := Client{}

	session, err := sess.NewAuthorisedSession(cl, api, auth)
	if err != nil {
		return nil, err
	}
	session = session.UpdateAToken(authtoken)
	client.Session = &session
	return &client, nil
}

func (cl Client) AddAPI(a api.API) Client {
	apis := cl.APIs
	apis[a.GetName()] = a
	cl.APIs = apis
	return cl
}
