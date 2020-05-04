package handles

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
	"time"
)

type (
	user struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
	timeSeries struct {
		Name string
		Time int64
		}
)

var (
	users = map[int]*user{}
	seq   = 1
)


func GetClaimUser(c echo.Context) (error, *timeSeries) {
	token := c.Get("user").(*jwt.Token)
	userClaim := token.Claims.(jwt.MapClaims)
	name := userClaim["name"].(string)
	exp := userClaim["exp"].(float64)
	return nil, &timeSeries{ Name: name, Time: int64(exp) }
}

func validateToken(c echo.Context) error {
	err, timeSeries := GetClaimUser(c)
	if err != nil {
		fmt.Println("Need to check with ERROR logs in ", err)
	}
	fmt.Println("Accessing with account is: ", timeSeries.Name)
	fmt.Println("Time to created Token by", timeSeries.Time)
	fmt.Println("Time to current by", time.Now().Unix())
	return nil
}

func CreateUser(c echo.Context) error {
	err := validateToken(c)
	if err != nil {
		fmt.Println("Failed to validate with Token")
	}

	u := &user{
		ID: seq,
	}
	if err := c.Bind(u); err != nil {
		return err
	}
	users[u.ID] = u
	seq++
	return c.JSON(http.StatusCreated, u)
}

func GetUser(c echo.Context) error {
	err := validateToken(c)
	if err != nil {
		fmt.Println("Failed to validate with Token")
	}
	id, _ := strconv.Atoi(c.Param("id"))
	fmt.Println("Hey: ", users[id])
	return c.JSON(http.StatusOK, users[id])
}

func UpdateUser(c echo.Context) error {
	u := new(user)
	if err := c.Bind(u); err != nil {
		return err
	}
	id, _ := strconv.Atoi(c.Param("id"))
	users[id].Name = u.Name
	return c.JSON(http.StatusOK, users[id])
}

func DeleteUser(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	delete(users, id)
	return c.NoContent(http.StatusNoContent)
}
