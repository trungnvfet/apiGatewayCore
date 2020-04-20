package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"time"

	"github.com/apiGatewayCore/handles"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)


func respondHeader(c echo.Context, u user) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	c.Response().WriteHeader(http.StatusOK)
	return json.NewEncoder(c.Response()).Encode(u)
}

func changedGenerator() string {
	out, err := exec.Command("uuidgen").Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("req-", out)
	output := "req-" + string(out)
	return output
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLength:  32,
		TokenLookup:  "header:" + echo.HeaderXCSRFToken,
		ContextKey:   "csrf",
		CookieName:   "_csrf",
		CookieMaxAge: 86400,
	}))
	e.Use(middleware.RequestIDWithConfig(middleware.RequestIDConfig{
		Generator: func() string {
			return changedGenerator()
		},
	}))

	e.POST("/users", handles.CreateUser)
	e.GET("/users/:id", handles.GetUser)
	e.PUT("/users/:id", handles.UpdateUser)
	e.DELETE("/users/:id", handles.DeleteUser)
	e.POST("/login/", handles.LoginHandler)
	e.HTTPErrorHandler = handles.CustomHTTPErrorHandler

	data, errs := json.MarshalIndent(e.Routes(), "", "  ")
	if errs != nil {
		return
	}
	ioutil.WriteFile("/tmp/routes.json", data, 0644)

	s := &http.Server{
		Addr:              ":2020",
		ReadTimeout:       30 * time.Minute,
		WriteTimeout:      30 * time.Minute,
		MaxHeaderBytes:    1000,
		IdleTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Minute,
	}
	e.Logger.Fatal(e.StartServer(s))
}
