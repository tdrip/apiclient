package sess

import (
	"fmt"
	"net/http"

	cfg "github.com/tdrip/apiclient/pkg/v1/cfg"
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
	if err != nil {
		return asess, err
	}
	asess.Session = sess
	asess.Session.GetHeaders = asess.AuthorisedHeaders
	asess.auth = auth
	return asess, err
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
	_, res, err := asess.Session.Get(asess.auth.Verifyauth)
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
	asess.Session.Get(asess.auth.Revokeauth)
}
