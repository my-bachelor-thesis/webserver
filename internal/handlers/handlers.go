package handlers

import (
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
	return c.Render(http.StatusOK, "login.html", newLoginData(c))
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

func Logout(c echo.Context) error {
	cookie := &http.Cookie{Name: "auth", Value: ""}
	c.SetCookie(cookie)
	return c.Redirect(http.StatusSeeOther, "/")
}

func RegisterGet(c echo.Context) error {
	return c.Render(http.StatusOK, "register.html", nil)
}

func RegisterPost(c echo.Context) error {
	user := &users.User{}
	if err := c.Bind(user); err != nil {
		return err
	}

	if user.Password != c.FormValue("password2") {
		return echo.NewHTTPError(http.StatusBadRequest, "passwords don't match")
	}

	encrypted, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(encrypted)

	if err = user.Insert(); err != nil {
		return err
	}

	return c.Redirect(http.StatusSeeOther, "/")
}

func IndexGet(c echo.Context) error {
	// load from all published from DB
	published, err := tasks.GetAllApprovedAndPublished()
	if err != nil {
		return err
	}
	data := newLoginData(c)
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
