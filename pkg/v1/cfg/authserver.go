package cfg

import (
	uris "github.com/tdrip/apiclient/pkg/v1/uris"
)

type AuthServer struct {
	EndPoint   uris.EndPoint
	Verifyauth string
	Revokeauth string
}
