package handlers

import (
	"github.com/labstack/echo/v4"
	"webserver/internal/jwt"
	"webserver/internal/postgres"
	"webserver/internal/postgres/rdg/users"
)

var emptySliceResponse = make([]int, 0)

type requestWithIdAndName struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func bindAndFind[T any](c echo.Context, getByIdFunc func(int) (T, error)) (*requestWithIdAndName, T, error) {
	var req requestWithIdAndName
	if err := c.Bind(&req); err != nil {
		return nil, *new(T), err
	}
	obj, err := getByIdFunc(req.Id)
	return &req, obj, err
}

func bindAndFindWithUserId[T any](c echo.Context, getByIdFunc func(int, int) (T, error), userId int) (*requestWithIdAndName, T, error) {
	var req requestWithIdAndName
	if err := c.Bind(&req); err != nil {
		return nil, *new(T), err
	}
	obj, err := getByIdFunc(req.Id, userId)
	return &req, obj, err
}

func getUserFromJWTCookie(tx postgres.PoolInterface, c echo.Context) (*users.User, error) {
	claims, err := jwt.GetClaimsFromRequest(c)
	if err != nil {
		return nil, err
	}
	return users.GetById(tx, claims.UserId)
}

func bindAndValidate[T any](c echo.Context, request T) error {
	if err := c.Bind(request); err != nil {
		return err
	}
	return c.Validate(request)
}
