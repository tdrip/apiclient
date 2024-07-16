package client

import (
	"context"
	"fmt"
	"testing"

	cfg "github.com/tdrip/apiclient/pkg/v1/cfg"
)

func TestClient(t *testing.T) {
	auth := cfg.AuthServer{}
	api, _ := cfg.NewAPIServer("jsonplaceholder.typicode.com", "todos")
	newclient := NewTlsSkip(api, auth)
	newclient.Session.Debug = true
	newclient.Session.DumpResponse = true
	newclient.Session.DumpRequest = true
	newclient.Session.Ctx = context.Background()
	bytes, resp, err := newclient.Session.Get("/1")
	if err != nil {
		t.Fatalf("%v", err.Error())
	}

	fmt.Println("Bytes")
	fmt.Printf("%s", string(bytes))
	fmt.Println(" ")

	fmt.Println("Response")
	fmt.Printf("%v", resp)
}
