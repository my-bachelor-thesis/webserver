package handlers

import (
	"github.com/labstack/echo/v4"
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
