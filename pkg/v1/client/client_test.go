package client

import (
	"fmt"
	"testing"

	cfg "github.com/tdrip/apiclient/pkg/v1/cfg"
)

func TestClient(t *testing.T) {
	auth := cfg.AuthServer{}
	api, _ := cfg.NewAPIServer("jsonplaceholder.typicode.com", "todos")
	newclient := NewTlsSkip(api, auth)
	newclient.AuthSession.Debug = true
	newclient.AuthSession.DumpResponse = true
	newclient.AuthSession.DumpRequest = true

	bytes, resp, err := newclient.AuthSession.Get("/1")
	if err != nil {
		t.Fatalf("%v", err.Error())
	}

	fmt.Println("Bytes")
	fmt.Printf("%s", string(bytes))
	fmt.Println(" ")

	fmt.Println("Response")
	fmt.Printf("%v", resp)
}
