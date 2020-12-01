package sonarcloud

import (
	"fmt"
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

func (sc *SonarClient) NewRequestWithParameters(method string, url string, params ...string) (*http.Request, error) {
	if l := len(params); l%2 != 0 {
		return nil, fmt.Errorf("params must be an even number, %d given", l)
	}

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Add("organization", sc.org)

	for i := 0; i < len(params); i++ {
		q.Add(params[i], params[i+1])
		i++
	}
	req.URL.RawQuery = q.Encode()

	req.SetBasicAuth(sc.token, "")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	return req, nil
}

func (sc *SonarClient) Do(req *http.Request) (*http.Response, error) {
	return sc.client.Do(req)
}
