package cfg

import (
	uris "github.com/tdrip/apiclient/pkg/v1/uris"
)

type APIServer struct {
	EndPoint uris.EndPoint
}

func NewAPIServer(server string, apipath string) (APIServer, error) {
	ai := APIServer{}
	api, err := uris.NewEndPoint(server, apipath)
	if err != nil {
		return ai, err
	}
	ai.EndPoint = api
	return ai, nil
}
