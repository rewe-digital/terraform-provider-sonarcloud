package sonarcloud

import (
	"fmt"
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

func (sc *SonarClient) NewRequest() (*http.Request, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/user_groups/search", API), nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Add("organization", sc.org)
	req.URL.RawQuery = q.Encode()

	req.SetBasicAuth(sc.token, "")

	return req, nil
}

func (sc *SonarClient) Do(req *http.Request) (*http.Response, error) {
	return sc.client.Do(req)
}
