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
	sess, err := NewAuthorisedSession(client, api, auth)
	sess.Session.Logger = logger
	return sess, err
}

func NewAuthorisedSession(client *http.Client, api cfg.APIServer, auth cfg.AuthServer) (AuthorisedSession, error) {
	asess := AuthorisedSession{}
	sess, err := NewSession(client, api)
	asess.Session = sess
	asess.auth = auth
	return asess, err
}

func (sess AuthorisedSession) PostBody(uri string, req interface{}) ([]byte, *http.Response, error) {
	return sess.AuthorizedRequest(http.MethodPost, uri, sess.Session.api.EndPoint, req)
}

func (sess AuthorisedSession) HeadBody(uri string, req interface{}) ([]byte, *http.Response, error) {
	return sess.AuthorizedRequest(http.MethodHead, uri, sess.Session.api.EndPoint, req)
}

func (sess AuthorisedSession) Get(uri string) ([]byte, *http.Response, error) {
	return sess.AuthorizedRequest(http.MethodGet, uri, sess.Session.api.EndPoint, nil)
}

func (sess AuthorisedSession) GetBody(uri string, req interface{}) ([]byte, *http.Response, error) {
	return sess.AuthorizedRequest(http.MethodGet, uri, sess.Session.api.EndPoint, req)
}

func (sess AuthorisedSession) PutBody(uri string, req interface{}) ([]byte, *http.Response, error) {
	return sess.AuthorizedRequest(http.MethodPut, uri, sess.Session.api.EndPoint, req)
}

func (sess AuthorisedSession) AuthorizedRequest(method string, uri string, ep uris.EndPoint, req interface{}) ([]byte, *http.Response, error) {

	url, err := ep.GetURL(uri)
	emptydata := []byte{}
	if err != nil {
		return emptydata, nil, err
	}
	return sess.Authorized(method, url, req)
}

func (sess AuthorisedSession) AuthorisedHeaders() map[string]string {

	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	headers["Authorization"] = "Bearer " + sess.accesstoken

	return headers
}

func (sess AuthorisedSession) Authorized(method string, url string, req interface{}) ([]byte, *http.Response, error) {
	emptydata := []byte{}

	res, err := sess.Session.APICall(method, url, req, sess.AuthorisedHeaders())

	if err != nil {
		return emptydata, res, err
	}

	if res == nil {
		return emptydata, res, fmt.Errorf("%s result was nil: %s Status Code %d", url, res.Status, res.StatusCode)
	}

	if sess.Session.Debug && sess.Session.Logger != nil {
		b, e := httputil.DumpResponse(res, sess.Session.DumpResponse)
		sess.Session.Logger("Authorized", b, e)
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

func (sess AuthorisedSession) UpdateAToken(accesstoken string) AuthorisedSession {
	sess.accesstoken = accesstoken
	return sess
}

func (sess AuthorisedSession) UpdateRToken(refreshtoken string) AuthorisedSession {
	sess.refreshtoken = refreshtoken
	return sess
}

func (sess AuthorisedSession) Verify() error {
	_, res, err := sess.AuthGet(sess.auth.Verifyauth)
	if err != nil {
		return err
	}

	if res.StatusCode == 200 {
		return nil
	} else {
		url, _ := sess.auth.EndPoint.GetURL("")
		return fmt.Errorf("%s - Response was not 200: %s Status Code %d", url, res.Status, res.StatusCode)
	}
}

func (sess AuthorisedSession) RevokeAndDisconnect() {
	sess.AuthGet(sess.auth.Revokeauth)
}

func (sess AuthorisedSession) AuthPostBody(uri string, req interface{}) ([]byte, *http.Response, error) {
	return sess.AuthorizedRequest(http.MethodPost, uri, sess.auth.EndPoint, req)
}

func (sess AuthorisedSession) AuthGet(uri string) ([]byte, *http.Response, error) {
	return sess.AuthorizedRequest(http.MethodGet, uri, sess.auth.EndPoint, nil)
}
