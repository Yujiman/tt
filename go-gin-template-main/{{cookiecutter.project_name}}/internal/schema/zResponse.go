package schema

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Empty = struct {}

//swagger:model http_response
type Response[T any] struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Result  T `json:"result"`
}

type ResponsePaginate[T any] struct {
	Status     bool     `json:"status"`
	StatusCode int      `json:"status_code"`
	Message    string   `json:"message"`
	Result     Paginate[T] `json:"result"`
}

type Paginate[T any] struct {
	Items T `json:"items"`
	Total int64       `json:"total"`
	Page  int         `json:"page"`
}

func Respond[T any](v T, c *gin.Context) error {
	c.JSON(http.StatusOK, Response[T]{
		Status:  true,
		Message: "Success",
		Result:  v,
	})

	return nil
}

func RespondPaginate[T any](v T, total int64, page int, c *gin.Context) {
	c.JSON(http.StatusOK, ResponsePaginate[T]{
		Status:     true,
		Message:    "OK",
		StatusCode: 0,
		Result: Paginate[T]{
			Items: v,
			Total: total,
			Page:  page,
		},
	})
}

func RespondMessage(v string, c *gin.Context) {
	c.JSON(http.StatusOK, Response[Empty]{
		Status:  true,
		Message: v,
	})
}
