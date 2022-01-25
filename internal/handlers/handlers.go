package handlers

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
	"webserver/internal/config"
	"webserver/internal/postgres/rdg/tasks"
	"webserver/internal/postgres/rdg/users"
)

func LoginGet(c echo.Context) error {
	return c.Render(http.StatusOK, "login.html", nil)
}

func LoginPost(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	// get user from DB
	user, err := users.GetByUsername(username)
	if err != nil {
		return err
	}

	// Throws unauthorized error
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return c.String(http.StatusBadRequest, "bad password")
	}

	// Set custom claims
	claims := &JwtCustomClaims{
		user.Id,
		user.IsAdmin,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(config.GetInstance().JWTSecret))
	if err != nil {
		return err
	}

	jwtCookie := &http.Cookie{}
	jwtCookie.Name = "auth"
	jwtCookie.Value = t
	jwtCookie.Expires = time.Now().Add(24 * time.Hour)
	c.SetCookie(jwtCookie)

	return c.Redirect(http.StatusSeeOther, "/")
}

func RegisterGet(c echo.Context) error {
	return c.Render(http.StatusOK, "register.html", nil)
}

func RegisterPost(c echo.Context) error {
	// insert user from DB
	return nil
}

func IndexGet(c echo.Context) error {
	// load from all published from DB
	published, err := tasks.GetAllApprovedAndPublished()
	if err != nil {
		return err
	}
	data := make(map[string]interface{})
	data["tasks"] = published
	return c.Render(http.StatusOK, "index.html", data)
}

func TaskGet(c echo.Context) error {
	return c.Render(http.StatusOK, "task.html", nil)
}

func AddTaskGet(c echo.Context) error {
	return c.Render(http.StatusOK, "add_task.html", nil)
}

func AddTaskPost(c echo.Context) error {
	// insert task into DB

	return c.Redirect(http.StatusSeeOther, "/")
}

// will get removed

func RestrictedGet(c echo.Context) error {
	claims := getClaimsFromRequest(c)
	return c.String(http.StatusOK, fmt.Sprintf("%#v", claims))
}
