package api

import (
	sess "github.com/tdrip/apiclient/pkg/session"
)

type API interface {
	GetSession() *sess.Session
	SetSession(sess *sess.Session)
	GetName() string
	SetName(s string)
}

type APIEndpoint struct {
	session *sess.Session
	Name    string
}

func NewAPIEndpoint(name string) APIEndpoint {
	api := APIEndpoint{}
	api.SetName(name)
	return api
}

func (aie APIEndpoint) HasSession() bool {
	return aie.session != nil
}

func (aie APIEndpoint) GetSession() *sess.Session {
	return aie.session
}

func (aie APIEndpoint) SetSession(sess *sess.Session) {
	aie.session = sess
}

func (aie APIEndpoint) GetName() string {
	return aie.Name
}

func (aie APIEndpoint) SetName(s string) {
	aie.Name = s
}
