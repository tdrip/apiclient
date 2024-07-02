package uris

import (
	"fmt"
	"net/url"
)

type EndPoint struct {
	uriPath string
	url     *url.URL
}

func NewEndPoint(server string, uriPath string) (EndPoint, error) {
	ep := EndPoint{}
	u, err := url.Parse(fmt.Sprintf("https://%s/", server))
	if err != nil {
		return ep, err
	}
	ep.url = u
	ep.uriPath = uriPath
	return ep, nil
}

func (ep EndPoint) GetURL(childuri string) (string, error) {

	newurl, newurlerr := ep.url.Parse(ep.uriPath + childuri)

	if newurlerr != nil {
		return "", fmt.Errorf("URL parsing error: Server %s, uri %s result %v", ep.url.String(), ep.uriPath+childuri, newurlerr)
	}
	return newurl.String(), nil
}
