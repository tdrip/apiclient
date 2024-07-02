package sess

import (
	"fmt"
	"net/http"
)

func (sess Session) UpdateAToken(accesstoken string) Session {
	sess.accesstoken = accesstoken
	return sess
}

func (sess Session) UpdateRToken(refreshtoken string) Session {
	sess.refreshtoken = refreshtoken
	return sess
}

func (sess *Session) Verify() error {
	_, res, err := sess.AuthGet(sess.verifyauth)
	if err != nil {
		return err
	}

	if res.StatusCode == 200 {
		return nil
	} else {
		url, _ := sess.auth.GetURL("")
		return fmt.Errorf("%s - Response was not 200: %s Status Code %d", url, res.Status, res.StatusCode)
	}
}

func (sess *Session) RevokeAndDisconnect() {
	sess.AuthGet(sess.revokeauth)
}

func (sess *Session) AuthPostBody(uri string, req interface{}) ([]byte, *http.Response, error) {
	return sess.AuthorizedRequest(http.MethodPost, uri, sess.auth, req)
}

func (sess *Session) AuthGet(uri string) ([]byte, *http.Response, error) {
	return sess.AuthorizedRequest(http.MethodGet, uri, sess.auth, nil)
}
