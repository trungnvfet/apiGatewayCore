package main

import (
	"encoding/json"
	"fmt"
	"github.com/apiGatewayCore/handles"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"time"
)

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
	//e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
	//	TokenLength:  32,
	//	TokenLookup:  "header:" + echo.HeaderXCSRFToken,
	//	ContextKey:   "csrf",
	//	CookieName:   "_csrf",
	//	//TokenLookup: "form:csrf",
	//	CookieMaxAge: 86400,
	//}))

	var IsValidate = middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte("secret"),
	})
	e.Use(middleware.RequestIDWithConfig(middleware.RequestIDConfig{
		Generator: func() string {
			return changedGenerator()
		},
	}))

	e.POST("/users", handles.CreateUser, IsValidate)
	e.POST("/login/", handles.LoginHandler)
	e.POST("/token/", handles.GenToken)
	e.GET("/users/:id", handles.GetUser, IsValidate)
	e.PUT("/users/:id", handles.UpdateUser)
	e.DELETE("/users/:id", handles.DeleteUser)
	e.HTTPErrorHandler = handles.CustomHTTPErrorHandler

	data, errs := json.MarshalIndent(e.Routes(), "", "  ")
	if errs != nil {
		return
	}
	err := ioutil.WriteFile("/tmp/routes.json", data, 0644)
	if err != nil {
		fmt.Println("Error to write json file in /tmp")
	}

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
