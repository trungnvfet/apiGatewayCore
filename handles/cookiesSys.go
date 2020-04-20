package handles

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo"
)

type (
	cookies struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	}
)

func readCookie(c echo.Context) (error, cookies) {
	cookie, err := c.Cookie("username")
	if err != nil {
		u := &cookies{
			Name:  "",
			Value: "",
		}
		return err, *u
	}
	fmt.Println(cookie.Name)
	fmt.Println(cookie.Value)
	u := &cookies{
		Name:  cookie.Name,
		Value: cookie.Value,
	}
	return nil, *u
}

func writeCookie(c echo.Context) error {
	cookie := new(http.Cookie)
	cookie.Name = "username"
	cookie.Value = "jon"
	cookie.Expires = time.Now().Add(24 * time.Hour)
	c.SetCookie(cookie)
	return c.String(http.StatusOK, "write a cookie")
}

func readAllCookies(c echo.Context) error {
	for _, cookie := range c.Cookies() {
		fmt.Println(cookie.Name)
		fmt.Println(cookie.Value)
	}
	return c.String(http.StatusOK, "read all the cookies")
}
