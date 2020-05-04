package handles

import (
	"net/http"

	"github.com/labstack/echo"
)

type errMsg struct {
	HostIP string
	Method string
	ErrorMsg string
	ErrorCode int
}

func CustomHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	msg := ""
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		msg = he.Message.(string)
	}
	c.Logger().Error(err)

	output := errMsg {
		c.RealIP(),
		c.Request().Method,
		msg,
		code,
	}

	if code == 400 {
		c.JSON(http.StatusBadRequest, output)
	}
	c.JSON(http.StatusInternalServerError, output)
}
