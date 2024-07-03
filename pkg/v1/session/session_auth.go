package sess

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"

	cfg "github.com/tdrip/apiclient/pkg/v1/cfg"
	uris "github.com/tdrip/apiclient/pkg/v1/uris"
)

type AuthorisedSession struct {
	Session
	accesstoken  string
	refreshtoken string
	auth         cfg.AuthServer
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
	return asess.Authorized(method, url, req)
}

func (asess AuthorisedSession) AuthorisedHeaders() map[string]string {

	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	headers["Authorization"] = "Bearer " + asess.accesstoken

	return headers
}

func (asess AuthorisedSession) Authorized(method string, url string, req interface{}) ([]byte, *http.Response, error) {
	emptydata := []byte{}

	res, err := asess.Session.APICall(method, url, req, asess.AuthorisedHeaders())

	if err != nil {
		return emptydata, res, err
	}

	if res == nil {
		return emptydata, res, fmt.Errorf("%s result was nil: %s Status Code %d", url, res.Status, res.StatusCode)
	}

	if asess.Session.Debug && asess.Session.Logger != nil {
		b, e := httputil.DumpResponse(res, asess.Session.DumpResponse)
		asess.Session.Logger("Authorized", b, e)
	}

	if res.StatusCode != 200 {
		return emptydata, res, fmt.Errorf("%s failed with Status: %s Status Code %d", url, res.Status, res.StatusCode)
	}

	if res.Body != http.NoBody {
		defer res.Body.Close()
		bytes, err := ioutil.ReadAll(res.Body)
		return bytes, res, err
	}
	// return empty if got no body, response and err
	return []byte{}, res, err
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
