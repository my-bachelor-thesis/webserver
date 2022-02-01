package handlers

import (
	"github.com/labstack/echo/v4"
	"webserver/internal/config"
)

func newLoginData(c echo.Context) map[string]interface{} {
	data := map[string]interface{}{}
	claims, err := getClaimsFromRequest(c)
	if err == nil {
		loginData := struct {
			IsAdmin bool
			Id      int
		}{IsAdmin: claims.IsAdmin, Id: claims.UserId}

		data["loginData"] = loginData
		data["isLoggedIn"] = true
	} else {
		data["isLoggedIn"] = false
	}
	return data
}

func getUserId(c echo.Context) (int, error) {
	if !config.GetInstance().IsProduction {
		return 1, nil
	}
	claims, err := getClaimsFromRequest(c)
	if err != nil {
		return 0, err
	}
	return claims.UserId, nil
}
