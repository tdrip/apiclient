package sess

import (
	"fmt"
	"net/http"

	cfg "github.com/tdrip/apiclient/pkg/v1/cfg"
	uris "github.com/tdrip/apiclient/pkg/v1/uris"
)

type AuthorisedSession struct {
	Session
	accesstoken    string
	refreshtoken   string
	auth           cfg.AuthServer
	GetAuthHeaders GetHeaders
}

func NewAuthorisedSessionCustomLogger(client *http.Client, api cfg.APIServer, auth cfg.AuthServer, logger SessionLog) (AuthorisedSession, error) {
	asess, err := NewAuthorisedSession(client, api, auth)
	asess.Session.Logger = logger
	return asess, err
}

func NewAuthorisedSession(client *http.Client, api cfg.APIServer, auth cfg.AuthServer) (AuthorisedSession, error) {
	asess := AuthorisedSession{}
	sess, err := NewSession(client, api)
	asess.Session = sess
	asess.auth = auth
	asess.GetAuthHeaders = asess.AuthorisedHeaders()
	return asess, err
}

func (asess AuthorisedSession) PostBody(uri string, req interface{}) ([]byte, *http.Response, error) {
	return asess.AuthorizedRequest(http.MethodPost, uri, asess.Session.api.EndPoint, req)
}

func (asess AuthorisedSession) HeadBody(uri string, req interface{}) ([]byte, *http.Response, error) {
	return asess.AuthorizedRequest(http.MethodHead, uri, asess.Session.api.EndPoint, req)
}

func (asess AuthorisedSession) Get(uri string) ([]byte, *http.Response, error) {
	return asess.AuthorizedRequest(http.MethodGet, uri, asess.Session.api.EndPoint, nil)
}

func (asess AuthorisedSession) GetBody(uri string, req interface{}) ([]byte, *http.Response, error) {
	return asess.AuthorizedRequest(http.MethodGet, uri, asess.Session.api.EndPoint, req)
}

func (asess AuthorisedSession) PutBody(uri string, req interface{}) ([]byte, *http.Response, error) {
	return asess.AuthorizedRequest(http.MethodPut, uri, asess.Session.api.EndPoint, req)
}

func (asess AuthorisedSession) AuthorizedRequest(method string, uri string, ep uris.EndPoint, req interface{}) ([]byte, *http.Response, error) {

	url, err := ep.GetURL(uri)
	emptydata := []byte{}
	if err != nil {
		return emptydata, nil, err
	}
	return asess.Session.Call(method, url, req, asess.GetAuthHeaders)
}

func (asess AuthorisedSession) AuthorisedHeaders() map[string]string {

	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	headers["Accept"] = "application/json"
	headers["Authorization"] = "Bearer " + asess.accesstoken

	return headers
}

func (asess AuthorisedSession) UpdateAToken(accesstoken string) AuthorisedSession {
	asess.accesstoken = accesstoken
	return asess
}

func (asess *AuthorisedSession) SetAToken(accesstoken string) {
	asess.accesstoken = accesstoken
}

func (asess AuthorisedSession) UpdateRToken(refreshtoken string) AuthorisedSession {
	asess.refreshtoken = refreshtoken
	return asess
}

func (asess AuthorisedSession) Verify() error {
	_, res, err := asess.AuthGet(asess.auth.Verifyauth)
	if err != nil {
		return err
	}

	if res.StatusCode == 200 {
		return nil
	} else {
		url, _ := asess.auth.EndPoint.GetURL("")
		return fmt.Errorf("%s - Response was not 200: %s Status Code %d", url, res.Status, res.StatusCode)
	}
}

func (asess AuthorisedSession) RevokeAndDisconnect() {
	asess.AuthGet(asess.auth.Revokeauth)
}

func (asess AuthorisedSession) AuthPostBody(uri string, req interface{}) ([]byte, *http.Response, error) {
	return asess.AuthorizedRequest(http.MethodPost, uri, asess.auth.EndPoint, req)
}

func (asess AuthorisedSession) AuthGet(uri string) ([]byte, *http.Response, error) {
	return asess.AuthorizedRequest(http.MethodGet, uri, asess.auth.EndPoint, nil)
}
