package sonarcloud

import (
	"io"
	"net/http"
	"time"
)

type SonarClient struct {
	client *http.Client
	org    string
	token  string
}

func NewSonarClient(org string, token string) (*SonarClient, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	return &SonarClient{
		client: client,
		org:    org,
		token:  token,
	}, nil
}

func (sc *SonarClient) NewRequest(method string, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Add("organization", sc.org)
	req.URL.RawQuery = q.Encode()

	req.SetBasicAuth(sc.token, "")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	return req, nil
}

func (sc *SonarClient) Do(req *http.Request) (*http.Response, error) {
	return sc.client.Do(req)
}
