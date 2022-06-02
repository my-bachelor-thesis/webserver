package handlers

import (
	"github.com/labstack/echo/v4"
)

var emptySliceResponse = make([]int, 0)

func bindAndValidate[T any](c echo.Context, request T) error {
	if err := c.Bind(request); err != nil {
		return err
	}
	return c.Validate(request)
}
