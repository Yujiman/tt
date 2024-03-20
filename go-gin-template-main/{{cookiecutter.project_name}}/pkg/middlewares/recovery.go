package middlewares

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func Recovery(handler func(c *gin.Context)) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					var se *os.SyscallError
					if errors.As(ne, &se) {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				if brokenPipe {
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) //nolint: errcheck
					c.Abort()
				} else {
					if e, ok := err.(error); ok {
						c.Error(e)
					} else {
						c.Error(fmt.Errorf("%s", err))
					}

					handler(c)
				}
			}
		}()

		c.Next()
	}
}
