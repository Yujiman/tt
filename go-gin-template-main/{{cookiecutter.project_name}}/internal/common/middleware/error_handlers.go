package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"git-ffd.kz/fmobile/ferr"
	"git-ffd.kz/pkg/goerr"
	"git-ffd.kz/pkg/gosentry"
	"github.com/gin-gonic/gin"
)

type ErrResponse struct {
	Status     bool        `json:"status"`
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Result     interface{} `json:"result"`
}

func (e ErrResponse) ToJson() (raw []byte, err error) {
	raw, err = json.Marshal(e)
	return raw, err
}

func GinErrorHandle(h func(c *gin.Context) error) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := h(c); err != nil {
			if len(c.Errors) == 0 {
				c.Error(err)
			}

			GinRecoveryFn(c)
		}
	}
}

func GinRecoveryFn(c *gin.Context) {
	err := c.Errors.Last()
	if err == nil {
		return
	}

	resp := &ErrResponse{
		Status:     false,
		StatusCode: 999999,
		Message:    err.Err.Error(),
		Result:     struct{}{},
	}

	// Если ошибка уже была обработана самим фреймворком, дописываем только тело ответа
	if c.IsAborted() {
		resp.StatusCode = ferr.ErrParseErrorBody.Code
		rawResp, _ := resp.ToJson()
		c.Writer.Write(rawResp)
		return
	}

	var goErr goerr.GoErr

	// Если ошибка это ошибка нашего приложения - возвращаем её с нужным статус-кодом
	if errors.As(err, &goErr) {
		if goErr.Code != 0 {
			resp.StatusCode = goErr.Code
		}

		if goErr.Result != nil {
			resp.Result = goErr.Result
		}

		if goErr.IsExpected {
			c.JSON(http.StatusOK, resp)
		} else {
			c.JSON(http.StatusOK, resp)
			// Все unexpected ошибки должны логироваться и попадать в Sentry
			internalServerErrorHandler(c.Request.Context(), err.Err, &gosentry.Info{
				Request:  c.Request,
				Response: c.Request.Response,
			})
		}

		return
	}

	// ErrorTypeBind - ошибка парсинга тела. Отдаем 400 код
	if err.Type == gin.ErrorTypeBind {
		resp.StatusCode = ferr.ErrParseErrorBody.Code
		c.JSON(http.StatusOK, resp)
		return
	}

	// В ином случае это неизвестная ошибка
	c.JSON(http.StatusOK, resp)
	internalServerErrorHandler(c.Request.Context(), err, &gosentry.Info{
		Request:  c.Request,
		Response: c.Request.Response,
	})
}

func internalServerErrorHandler(ctx context.Context, err error, sentryInfo *gosentry.Info) {
	gosentry.SendError(ctx, err, sentryInfo)
}
