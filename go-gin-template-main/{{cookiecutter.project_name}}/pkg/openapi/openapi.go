package openapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type OpenApi struct {
	client          *http.Client
	err             error
	openApiEndpoint string
	AppName         string                 `json:"service_name"`
	Doc             map[string]interface{} `json:"doc"`
}

type Response struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Result  interface{} `json:"result"`
}

type IOpenApi interface {
	Send(ctx context.Context) IOpenApi
	Error() error
}

func NewOpenApiClient(appName, openApiEndpoint, doc string) IOpenApi {
	buff := make(map[string]interface{})
	err := json.Unmarshal([]byte(doc), &buff)
	return &OpenApi{
		client: &http.Client{
			Timeout: time.Second * 30,
		},
		AppName:         appName,
		openApiEndpoint: openApiEndpoint,
		Doc:             buff,
		err:             err,
	}
}

func (s OpenApi) Send(ctx context.Context) IOpenApi {
	if s.err != nil {
		return s
	}
	var resBody Response
	req, err := s.newRequest(ctx, http.MethodPost, s.openApiEndpoint, s)
	if err != nil {
		s.err = err
		return s
	}
	resp, err := s.client.Do(req)
	if err != nil {
		s.err = err
		return s
	}
	if err = json.NewDecoder(resp.Body).Decode(&resBody); err != nil {
		s.err = err
		return s
	}
	if !resBody.Status || resp.StatusCode != http.StatusOK {
		s.err = fmt.Errorf(resBody.Message)
		return s
	}

	return s
}

func (s OpenApi) Error() error {
	if s.err != nil {
		return s.err
	}
	return nil
}

func (s OpenApi) newRequest(ctx context.Context, method string, url string, body interface{}) (*http.Request, error) {
	var bodyReader io.Reader
	if body != nil {
		rawBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewBuffer(rawBody)
	}
	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, err
	}
	return req, nil
}
