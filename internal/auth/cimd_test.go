package auth

import (
	"context"
	"io"
	"net"
	"net/http"
	"strings"
	"testing"
)

func TestValidateCIMDURLRejectsUnsafeTargets(t *testing.T) {
	lookup := func(context.Context, string) ([]net.IP, error) {
		return []net.IP{net.ParseIP("203.0.113.10")}, nil
	}
	cases := []string{
		"http://client.example/app.json",
		"https://client.example",
		"https://client.example/",
		"https://127.0.0.1/app.json",
		"https://169.254.169.254/latest",
		"https://10.1.1.1/app.json",
		"https://user:pass@client.example/app.json",
	}
	for _, raw := range cases {
		if err := ValidateCIMDURL(raw, lookup); err == nil {
			t.Fatalf("%s foi aceito", raw)
		}
	}
	if err := ValidateCIMDURL("https://client.example/app.json", lookup); err != nil {
		t.Fatal(err)
	}
	private := func(context.Context, string) ([]net.IP, error) {
		return []net.IP{net.ParseIP("192.168.0.5")}, nil
	}
	if err := ValidateCIMDURL("https://client.example/app.json", private); err == nil {
		t.Fatal("DNS privado foi aceito")
	}
}

func TestFetchRejectsMismatchedClientID(t *testing.T) {
	f := CIMDFetcher{
		Lookup: func(context.Context, string) ([]net.IP, error) {
			return []net.IP{net.ParseIP("203.0.113.10")}, nil
		},
		Client: &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
			body := `{"client_id":"https://evil.example/app.json","client_name":"x","redirect_uris":["http://127.0.0.1/cb"]}`
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     http.Header{"Content-Type": []string{"application/json"}},
				Body:       io.NopCloser(strings.NewReader(body)),
			}, nil
		})},
	}
	if _, _, err := f.Fetch(context.Background(), "https://client.example/app.json"); err == nil {
		t.Fatal("client_id divergente foi aceito")
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }
