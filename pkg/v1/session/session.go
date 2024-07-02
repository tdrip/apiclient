package sess

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"

	cfg "github.com/tdrip/apiclient/pkg/v1/cfg"
	uris "github.com/tdrip/apiclient/pkg/v1/uris"
)

type SessionLog func(msg string, data []byte, err error)

type Session struct {
	accesstoken  string
	refreshtoken string
	api          cfg.APIServer
	client       *http.Client
	Debug        bool
	DumpResponse bool
	DumpRequest  bool

	Logger SessionLog
}

func NewSessionCustomLogger(client *http.Client, api cfg.APIServer, logger SessionLog) (Session, error) {
	sess, err := NewSession(client, api)
	sess.Logger = logger
	return sess, err
}

func NewSession(client *http.Client, api cfg.APIServer) (Session, error) {
	sess := Session{}
	sess.client = client
	sess.api = api
	return sess, nil
}

func (sess Session) PostBody(uri string, req interface{}) ([]byte, *http.Response, error) {
	return sess.UnAuthorizedRequest(http.MethodPost, uri, sess.api.EndPoint, req)
}

func (sess Session) HeadBody(uri string, req interface{}) ([]byte, *http.Response, error) {
	return sess.UnAuthorizedRequest(http.MethodHead, uri, sess.api.EndPoint, req)
}

func (sess Session) Get(uri string) ([]byte, *http.Response, error) {
	return sess.UnAuthorizedRequest(http.MethodGet, uri, sess.api.EndPoint, nil)
}

func (sess Session) GetBody(uri string, req interface{}) ([]byte, *http.Response, error) {
	return sess.UnAuthorizedRequest(http.MethodGet, uri, sess.api.EndPoint, req)
}

func (sess Session) PutBody(uri string, req interface{}) ([]byte, *http.Response, error) {
	return sess.UnAuthorizedRequest(http.MethodPut, uri, sess.api.EndPoint, req)
}

func (sess Session) UnAuthorizedRequest(method string, uri string, ep uris.EndPoint, req interface{}) ([]byte, *http.Response, error) {

	url, err := ep.GetURL(uri)
	emptydata := []byte{}
	if err != nil {
		return emptydata, nil, err
	}
	return sess.UnAuthorized(method, url, req)
}

func (sess Session) UnAuthorisedHeaders() map[string]string {
	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	return headers
}

func (sess Session) UnAuthorized(method string, url string, req interface{}) ([]byte, *http.Response, error) {
	emptydata := []byte{}

	res, err := sess.APICall(method, url, req, sess.UnAuthorisedHeaders())

	if err != nil {
		return emptydata, res, err
	}

	if res == nil {
		return emptydata, res, fmt.Errorf("%s result was nil: %s Status Code %d", url, res.Status, res.StatusCode)
	}

	if sess.Debug && sess.Logger != nil {
		b, e := httputil.DumpResponse(res, sess.DumpResponse)
		sess.Logger("UnAuthorized", b, e)
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

func (sess Session) APICall(method string, url string, body interface{}, headers map[string]string) (*http.Response, error) {

	cs, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	tosent := string(cs)
	payload := strings.NewReader(tosent)
	req, reqerr := http.NewRequest(method, url, payload)

	if reqerr != nil {
		return nil, reqerr
	}

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	if sess.Debug && sess.Logger != nil {
		b, e := httputil.DumpRequestOut(req, sess.DumpRequest)
		sess.Logger("APICall", b, e)

	}

	if sess.client == nil {
		return http.DefaultClient.Do(req)
	} else {
		return sess.client.Do(req)
	}
}

func DefaultLogger(msg string, data []byte, err error) {
	if err == nil {
		log.Printf("Message (Error) %s\n  Data:  %s\n   Err:  %s\n", msg, data, err.Error())
	} else {
		log.Printf("Message (Debug):%s\n  Data:  %s\n", msg, data)
	}
}
