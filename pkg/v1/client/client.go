package client

import (
	"crypto/tls"
	"net/http"
	"time"

	api "github.com/tdrip/apiclient/pkg/v1/api"
	sess "github.com/tdrip/apiclient/pkg/v1/session"
)

type Client struct {
	Session *sess.Session
	APIs    map[string]api.API
}

func buildClient() *http.Client {
	transport := &http.Transport{Proxy: http.ProxyFromEnvironment, TLSClientConfig: &tls.Config{Renegotiation: tls.RenegotiateOnceAsClient, InsecureSkipVerify: true}, TLSNextProto: map[string]func(authority string, c *tls.Conn) http.RoundTripper{}}
	return &http.Client{
		Timeout:   time.Second * 10,
		Transport: transport,
	}

}

func NewTlsSkip(server string, authserver string, authtoken string, apipath string, authpath string, verifyauth string, revokeauth string) (*Client, error) {
	return New(buildClient(), server, authserver, authtoken, apipath, authpath, verifyauth, revokeauth)
}

func New(cl *http.Client, server string, authserver string, authtoken string, apipath string, authpath string, verifyauth string, revokeauth string) (*Client, error) {
	client := Client{}
	session, err := sess.NewSession(cl, server, authserver, apipath, authpath, verifyauth, revokeauth)
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
