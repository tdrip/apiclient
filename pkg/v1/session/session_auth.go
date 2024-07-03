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

func NewAuthorisedSessionCustomLogger(client *http.Client, api cfg.APIServer, auth cfg.AuthServer, logger SessionLog) AuthorisedSession {
	asess := NewAuthorisedSession(client, api, auth)
	asess.Logger = logger
	return asess
}

func NewAuthorisedSession(client *http.Client, api cfg.APIServer, auth cfg.AuthServer) AuthorisedSession {
	asess := AuthorisedSession{}
	asess.Session = NewSession(client, api)
	asess.GetHeaders = asess.AuthorisedHeaders
	asess.auth = auth
	return asess
}

func (asess AuthorisedSession) AuthorisedHeaders() map[string]string {

	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	headers["Accept"] = "application/json"
	if len(asess.accesstoken) > 0 {
		headers["Authorization"] = "Bearer " + asess.accesstoken
	}
	if asess.Session.Debug && asess.Session.Logger != nil {
		for k, v := range headers {
			asess.Session.Logger("AuthorisedHeaders", fmt.Sprintf("[%s] : %s", k, v), nil)
		}
	}
	return headers
}

func (asess AuthorisedSession) UpdateAToken(accesstoken string) AuthorisedSession {
	asess.accesstoken = accesstoken
	return asess
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

func SetAToken(asess *AuthorisedSession, accesstoken string) {
	asess.accesstoken = accesstoken
}
