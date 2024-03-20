// Пример написания клиента для интеграции с внешними сервисами

package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"git-ffd.kz/fmobile/ferr"
	"git-ffd.kz/pkg/clientrip"
	"git-ffd.kz/pkg/goerr"
	"git-ffd.kz/pkg/golog"
	"git-ffd.kz/pkg/golog/contrib/clientriplog"
	"git-ffd.kz/pkg/requestid"

	"{{cookiecutter.package_name}}/internal/schema/onec"
)

type OneC interface {
	Reservation(ctx context.Context, data onec.ReservationRequest) (onec.ReservationResponse, error)
	GetOrderStatus(ctx context.Context, data onec.GetOrderStatusRequest) (onec.OrderStatusResponse, error)
}

type OneCClient struct {
	client  *http.Client
	baseURL *url.URL
	token   string
}

func NewOneCClient(rawBaseURL string, token string, logger golog.ContextLogger) (*OneCClient, error) {
	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		return nil, fmt.Errorf("url.Parse(%q): %w", rawBaseURL, err)
	}

	return &OneCClient{
		client: &http.Client{
			Transport: clientrip.Chain(
				http.DefaultTransport,
				requestid.ClientMiddleware(),
				clientriplog.LoggingMiddleware(logger),
			),
			Timeout: time.Second * 60,
		},
		baseURL: baseURL,
		token:   token,
	}, nil
}

func (c *OneCClient) Reservation(ctx context.Context, data onec.ReservationRequest) (onec.ReservationResponse, error) {
	req, err := c.newRequest(ctx, http.MethodPost, c.baseURL.JoinPath("/api/v1/reservation"), data)
	if err != nil {
		return onec.ReservationResponse{}, ferr.ErrOneCIntegrationError.WithErr(err).WithCtx(ctx).
			WithMsgf("newRequest")
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return onec.ReservationResponse{}, ferr.ErrOneCIntegrationError.WithErr(err).WithCtx(ctx).
			WithMsgf("client.Do")
	}

	defer resp.Body.Close()

	rawResponse, err := io.ReadAll(resp.Body)
	if err != nil {
		return onec.ReservationResponse{}, ferr.ErrOneCIntegrationError.WithErr(err).WithResp(resp).
			WithMsgf("read body")
	}

	if resp.StatusCode != http.StatusOK {
		return onec.ReservationResponse{}, ferr.ErrTerminalIntegration.WithErr(err).WithResp(resp).
			WithResult(map[string]interface{}{"response": string(rawResponse)}).
			WithMsgf("unexpected http status code")
	}

	var response onec.ReservationResponse
	if err := json.Unmarshal(rawResponse, &response); err != nil {
		return onec.ReservationResponse{}, goerr.Wrap(err).WithErr(err).WithResp(resp).
			WithResult(map[string]interface{}{"response": string(rawResponse)}).
			WithMsgf("unmarshal error")
	}

	return response, nil
}

func (c *OneCClient) GetOrderStatus(
	ctx context.Context,
	data onec.GetOrderStatusRequest,
) (onec.OrderStatusResponse, error) {
	query := make(url.Values, 1)
	query.Set("order_id", strconv.Itoa(int(data.OrderID)))

	requestURL := c.baseURL.JoinPath("/api/v1/order-status")
	requestURL.RawQuery = query.Encode()

	req, err := c.newRequest(ctx, http.MethodGet, requestURL, data)
	if err != nil {
		return onec.OrderStatusResponse{}, ferr.ErrOneCIntegrationError.WithErr(err).WithCtx(ctx).
			WithMsgf("newRequest")
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return onec.OrderStatusResponse{}, ferr.ErrOneCIntegrationError.WithErr(err).WithCtx(ctx).
			WithMsgf("client.Do")
	}

	defer resp.Body.Close()

	rawResponse, err := io.ReadAll(resp.Body)
	if err != nil {
		return onec.OrderStatusResponse{}, ferr.ErrOneCIntegrationError.WithErr(err).WithResp(resp).
			WithMsgf("read body")
	}

	if resp.StatusCode != http.StatusOK {
		return onec.OrderStatusResponse{}, ferr.ErrTerminalIntegration.WithErr(err).WithResp(resp).
			WithResult(map[string]interface{}{"response": string(rawResponse)}).
			WithMsgf("unexpected http status code")
	}

	var response onec.OrderStatusResponse
	if err := json.Unmarshal(rawResponse, &response); err != nil {
		return onec.OrderStatusResponse{}, goerr.Wrap(err).WithErr(err).WithResp(resp).
			WithResult(map[string]interface{}{"response": string(rawResponse)}).
			WithMsgf("unmarshal error")
	}

	return response, nil
}

func (c *OneCClient) newRequest(
	ctx context.Context,
	method string,
	requestUrl *url.URL,
	body interface{},
) (*http.Request, error) {
	var bodyReader io.Reader
	if body != nil {
		rawBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}

		bodyReader = bytes.NewBuffer(rawBody)
	}

	req, err := http.NewRequestWithContext(ctx, method, requestUrl.String(), bodyReader)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Add("Content-Type", "application/json")
	}

	req.Header.Set("Authorization", c.token)

	return req, nil
}
