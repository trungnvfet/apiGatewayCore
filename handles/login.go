package handles

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo"
	"io/ioutil"
	"log"
	"net/http"
)

type loginSchema struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func loginUser(username string, password string) (bool, error) {
	fmt.Println(username + password)
	return true, nil
}

func LoginHandler(c echo.Context) error {
	w := c.Response()
	r := c.Request()
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return nil
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Body read error, %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return nil
	}
	var schema loginSchema
	if err = json.Unmarshal(body, &schema); err != nil {
		log.Printf("Body parse error, %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return nil
	}

	ok, err := loginUser(schema.Username, schema.Password)
	if err != nil {
		log.Printf("Login user DB error, %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return nil
	}

	if !ok {
		log.Printf("Unauthorized access for user: %v", schema.Username)
		w.WriteHeader(http.StatusUnauthorized)
		return nil
	}
	w.WriteHeader(http.StatusOK)
	return nil
}
