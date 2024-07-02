package cfg

import (
	uris "github.com/tdrip/apiclient/pkg/v1/uris"
)

type AuthServer struct {
	EndPoint   uris.EndPoint
	Verifyauth string
	Revokeauth string
}

func NewAuthServer(server string, authpath string) (AuthServer, error) {
	ai := AuthServer{}
	api, err := uris.NewEndPoint(server, authpath)
	if err != nil {
		return ai, err
	}
	ai.EndPoint = api
	return ai, nil
}

func NewAuthServerWithVerify(server string, authpath string, verifyauth string, revokeauth string) (AuthServer, error) {
	ai := AuthServer{}
	api, err := uris.NewEndPoint(server, authpath)
	if err != nil {
		return ai, err
	}
	ai.EndPoint = api
	ai.Revokeauth = revokeauth
	ai.Verifyauth = verifyauth
	return ai, nil
}
