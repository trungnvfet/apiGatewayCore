package handles

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

type jwtCustomClaims struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	jwt.StandardClaims
}

func GenToken(c echo.Context) error {
	// Validate username & password with input
	username := c.FormValue("username")
	password := c.FormValue("password")

	if username != "admin" || password != "Viettelidc@2020" {
		panic("Failed to put your admin account")
	}

	if username == password {
		panic("Failed to put same username and password together")
	}

	if len(username) == 0 || len(password) == 0 {
		panic("Failed to put your admin account")
	}

	claims := &jwtCustomClaims{
		"Admin",
		true,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenIssue, err := token.SignedString([]byte("secret"))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, echo.Map{
		"token": "Bearer " + tokenIssue,
	})
}