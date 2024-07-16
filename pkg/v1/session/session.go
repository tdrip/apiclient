package sess

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"

	utils "github.com/tdrip/apiclient/pkg/v1/utils"

	cfg "github.com/tdrip/apiclient/pkg/v1/cfg"
	uris "github.com/tdrip/apiclient/pkg/v1/uris"
)

type SessionLog func(msg string, data string, err error)
type GetHeaders func(sess Session) map[string]string
type MakeRequest func(sess Session, method string, uri string, ep uris.EndPoint, req interface{}) ([]byte, *http.Response, error)

type Session struct {
	api  cfg.APIServer
	auth cfg.AuthServer

	Ctx                    context.Context
	ContextAccessTokenName string

	client       *http.Client
	Debug        bool
	DumpResponse bool
	DumpRequest  bool

	Logger      SessionLog
	GetHeaders  GetHeaders
	MakeRequest MakeRequest

	accesstoken  string
	refreshtoken string
}

func NewSessionCustomLogger(client *http.Client, api cfg.APIServer, auth cfg.AuthServer, logger SessionLog) Session {
	sess := NewSession(client, api, auth)
	sess.Logger = logger
	return sess
}

func NewSession(client *http.Client, api cfg.APIServer, auth cfg.AuthServer) Session {
	sess := Session{}
	sess.client = client
	sess.api = api
	sess.auth = auth
	sess.GetHeaders = DefaultHeaders
	sess.MakeRequest = DoRequest
	return sess
}

func (sess Session) PostBody(uri string, req interface{}) ([]byte, *http.Response, error) {
	return sess.MakeRequest(sess, http.MethodPost, uri, sess.api.EndPoint, req)
}

func (sess Session) HeadBody(uri string, req interface{}) ([]byte, *http.Response, error) {
	return sess.MakeRequest(sess, http.MethodHead, uri, sess.api.EndPoint, req)
}

func (sess Session) Get(uri string) ([]byte, *http.Response, error) {
	return sess.MakeRequest(sess, http.MethodGet, uri, sess.api.EndPoint, nil)
}

func (sess Session) GetBody(uri string, req interface{}) ([]byte, *http.Response, error) {
	return sess.MakeRequest(sess, http.MethodGet, uri, sess.api.EndPoint, req)
}

func (sess Session) PutBody(uri string, req interface{}) ([]byte, *http.Response, error) {
	return sess.MakeRequest(sess, http.MethodPut, uri, sess.api.EndPoint, req)
}

func DoRequest(sess Session, method string, uri string, ep uris.EndPoint, req interface{}) ([]byte, *http.Response, error) {
	url, err := ep.GetURL(uri)
	emptydata := []byte{}
	if err != nil {
		return emptydata, nil, err
	}
	return sess.Call(method, url, req, sess.GetHeaders(sess))
}

func DefaultHeaders(sess Session) map[string]string {
	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	headers["Accept"] = "application/json"
	if len(sess.accesstoken) > 0 {
		headers["Authorization"] = "Bearer " + sess.accesstoken
	}
	if sess.Ctx != nil {
		if auth, ok := sess.Ctx.Value(sess.ContextAccessTokenName).(string); ok {
			headers["Authorization"] = "Bearer " + auth
		}
	}
	if sess.Debug && sess.Logger != nil {
		for k, v := range headers {
			sess.Logger("DefaultHeaders", fmt.Sprintf("[%s] : %s", k, v), nil)
		}
	}

	return headers
}

func (sess Session) Call(method string, url string, req interface{}, headers map[string]string) ([]byte, *http.Response, error) {
	emptydata := []byte{}

	res, err := sess.APICall(method, url, req, headers)

	if res != nil && sess.Debug && sess.Logger != nil {
		b, e := httputil.DumpResponse(res, sess.DumpResponse)
		sess.Logger("Call", string(b), e)
	}

	if res != nil && res.Body != http.NoBody {
		defer res.Body.Close()
		bytes, err := io.ReadAll(res.Body)
		return bytes, res, err
	}

	// return empty if got no body, response and err
	return emptydata, res, err
}

func (sess Session) APICall(method string, url string, body interface{}, headers map[string]string) (*http.Response, error) {

	req, reqerr := http.NewRequest(method, url, nil)

	if reqerr != nil {
		return nil, reqerr
	}

	if body != nil {
		ct := headers["Content-Type"]
		if strings.Contains(ct, "json") {
			cs, err := json.Marshal(body)
			if err != nil {
				return nil, err
			}

			tosent := string(cs)
			payload := strings.NewReader(tosent)
			req, reqerr = http.NewRequest(method, url, payload)

			if reqerr != nil {
				return nil, reqerr
			}
		}
	}

	if sess.Ctx != nil {
		// add context to the request
		req = req.WithContext(sess.Ctx)
	}

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	if sess.Debug && sess.Logger != nil {
		b, e := httputil.DumpRequestOut(req, sess.DumpRequest)
		sess.Logger("APICall", string(b), e)
	}

	if sess.client == nil {
		return http.DefaultClient.Do(req)
	} else {
		return sess.client.Do(req)
	}

}

func DefaultLogger(msg string, data string, err error) {
	if err != nil {
		log.Printf("Message (Error) %s\n  Data:  %s\n   Err:  %s\n", msg, data, err.Error())
	} else {
		log.Printf("Message (Debug):%s\n  Data:  %s\n", msg, data)
	}
}

func (sess Session) HasToken() bool {
	return len(sess.accesstoken) > 0
}

func (sess Session) UpdateAToken(accesstoken string) Session {
	sess.accesstoken = accesstoken
	return sess
}

func (sess Session) UpdateRToken(refreshtoken string) Session {
	sess.refreshtoken = refreshtoken
	return sess
}

func (sess Session) Verify() error {
	_, res, err := sess.Get(sess.auth.Verifyauth)
	if err != nil {
		return err
	}

	if !utils.RequestIsSuccessful(res.StatusCode) {
		return nil
	} else {
		url, _ := sess.auth.EndPoint.GetURL("")
		return fmt.Errorf("%s - Response was not a 2xx success code: %s Status Code %d", url, res.Status, res.StatusCode)
	}
}

func (sess Session) RevokeAndDisconnect() {
	sess.Get(sess.auth.Revokeauth)
}
