// Пример интеграции с нашими внутренними сервисами

package repository

//import (
//	"bytes"
//	"context"
//	"encoding/json"
//	"fmt"
//	"io"
//	"net/http"
//	"net/url"
//	"time"
//
//	"git-ffd.kz/fmobile/ferr"
//	"git-ffd.kz/pkg/clientrip"
//	"git-ffd.kz/pkg/goauth"
//	"git-ffd.kz/pkg/goerr"
//	"git-ffd.kz/pkg/golog"
//	"git-ffd.kz/pkg/golog/contrib/clientriplog"
//	"git-ffd.kz/pkg/requestid"
//
//	"{{cookiecutter.package_name}}/internal/schema/onlinepayment"
//)
//
//type OnlinePayment interface {
//	CheckPayment(ctx context.Context, data onlinepayment.CheckPaymentRequest) (onlinepayment.CheckPaymentResult, error)
//}
//
//type OnlinePaymentClient struct {
//	client  *http.Client
//	baseURL *url.URL
//}
//
//func NewOnlinePaymentClient(rawBaseURL string, logger golog.ContextLogger) (*OnlinePaymentClient, error) {
//	baseURL, err := url.Parse(rawBaseURL)
//	if err != nil {
//		return nil, fmt.Errorf("url.Parse(%q): %w", rawBaseURL, err)
//	}
//
//	return &OnlinePaymentClient{
//		client: &http.Client{
//			Transport: clientrip.Chain(
//				http.DefaultTransport,
//				requestid.ClientMiddleware(),
//				clientriplog.LoggingMiddleware(logger),
//			),
//			Timeout: time.Second * 60,
//		},
//		baseURL: baseURL,
//	}, nil
//}
//
//func (c *OnlinePaymentClient) CheckPayment(
//	ctx context.Context,
//	data onlinepayment.CheckPaymentRequest,
//) (onlinepayment.CheckPaymentResult, error) {
//	req, err := c.newRequest(ctx, http.MethodPost, c.baseURL.JoinPath("/api/v1/payments/check"), data)
//	if err != nil {
//		return onlinepayment.CheckPaymentResult{}, ferr.ErrOnlineInterationError.WithErr(err).WithCtx(ctx).
//			WithMsgf("newRequest")
//	}
//
//	resp, err := c.client.Do(req)
//	if err != nil {
//		return onlinepayment.CheckPaymentResult{}, ferr.ErrOnlineInterationError.WithErr(err).WithCtx(ctx).
//			WithMsgf("client.Do")
//	}
//
//	defer resp.Body.Close()
//
//	rawResponse, err := io.ReadAll(resp.Body)
//	if err != nil {
//		return onlinepayment.CheckPaymentResult{}, ferr.ErrOnlineInterationError.WithErr(err).WithResp(resp).
//			WithMsgf("read body")
//	}
//
//	if resp.StatusCode != http.StatusOK {
//		return onlinepayment.CheckPaymentResult{}, ferr.ErrOnlineInterationError.WithErr(err).WithResp(resp).
//			WithResult(map[string]interface{}{"response": string(rawResponse)}).
//			WithMsgf("unexpected http status code")
//	}
//
//	if err = goerr.RaiseForStatus(rawResponse); err != nil {
//		return onlinepayment.CheckPaymentResult{}, goerr.Wrap(err).WithCtx(ctx).WithResp(resp)
//	}
//
//	var response onlinepayment.Response[onlinepayment.CheckPaymentResult]
//	if err := json.Unmarshal(rawResponse, &response); err != nil {
//		return onlinepayment.CheckPaymentResult{}, ferr.ErrOnlineInterationError.WithErr(err).WithResp(resp).
//			WithResult(map[string]interface{}{"response": string(rawResponse)}).
//			WithMsgf("unmarshal error")
//	}
//
//	return response.Result, nil
//}
//
//func (c *OnlinePaymentClient) newRequest(
//	ctx context.Context,
//	method string,
//	requestUrl *url.URL,
//	body interface{},
//) (*http.Request, error) {
//	var bodyReader io.Reader
//	if body != nil {
//		rawBody, err := json.Marshal(body)
//		if err != nil {
//			return nil, err
//		}
//
//		bodyReader = bytes.NewBuffer(rawBody)
//	}
//
//	req, err := http.NewRequestWithContext(ctx, method, requestUrl.String(), bodyReader)
//	if err != nil {
//		return nil, err
//	}
//
//	if body != nil {
//		req.Header.Add("Content-Type", "application/json")
//	}
//
//	authHeader := goauth.AuthHeaderFromContext(ctx)
//	if authHeader != "" {
//		req.Header.Add("Authorization", authHeader)
//	}
//
//	return req, nil
//}
