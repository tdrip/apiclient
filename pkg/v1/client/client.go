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

type Session struct {
	accesstoken  string
	refreshtoken string
	api          cfg.APIServer
	auth         cfg.AuthServer
	client       *http.Client
	Debug        bool
	DumpResponse bool
	DumpRequest  bool
}

func NewSession(client *http.Client, api cfg.APIServer, auth cfg.AuthServer) (Session, error) {
	sess := Session{}
	sess.client = client
	sess.api = api
	sess.auth = auth
	return sess, nil
}

func (sess Session) PostBody(uri string, req interface{}) ([]byte, *http.Response, error) {
	return sess.AuthorizedRequest(http.MethodPost, uri, sess.api.EndPoint, req)
}

func (sess Session) HeadBody(uri string, req interface{}) ([]byte, *http.Response, error) {
	return sess.AuthorizedRequest(http.MethodHead, uri, sess.api.EndPoint, req)
}

func (sess Session) Get(uri string) ([]byte, *http.Response, error) {
	return sess.AuthorizedRequest(http.MethodGet, uri, sess.api.EndPoint, nil)
}

func (sess Session) GetBody(uri string, req interface{}) ([]byte, *http.Response, error) {
	return sess.AuthorizedRequest(http.MethodGet, uri, sess.api.EndPoint, req)
}

func (sess Session) PutBody(uri string, req interface{}) ([]byte, *http.Response, error) {
	return sess.AuthorizedRequest(http.MethodPut, uri, sess.api.EndPoint, req)
}

func (sess Session) AuthorizedRequest(method string, uri string, ep uris.EndPoint, req interface{}) ([]byte, *http.Response, error) {

	url, err := ep.GetURL(uri)
	emptydata := []byte{}
	if err != nil {
		return emptydata, nil, err
	}
	return sess.Authorized(method, url, req)
}

func (sess Session) Authorized(method string, url string, req interface{}) ([]byte, *http.Response, error) {
	emptydata := []byte{}

	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	headers["Authorization"] = "Bearer " + sess.accesstoken

	res, err := sess.APICall(method, url, req, headers)

	if err != nil {
		return emptydata, res, err
	}

	if res == nil {
		return emptydata, res, fmt.Errorf("%s result was nil: %s Status Code %d", url, res.Status, res.StatusCode)
	}

	if sess.Debug {
		debug(httputil.DumpResponse(res, sess.DumpResponse))
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

	if sess.Debug {
		debug(httputil.DumpRequestOut(req, sess.DumpRequest))
	}

	if sess.client == nil {
		return http.DefaultClient.Do(req)
	} else {
		return sess.client.Do(req)
	}
}

func debug(data []byte, err error) {
	if err == nil {
		log.Printf("%s\n\n", data)
	} else {
		log.Fatalf("%s\n\n", err)
	}
}
