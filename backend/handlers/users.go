package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *Handlers) CreateUser(c echo.Context) error {
	return c.String(http.StatusOK, "Not implemented")
}

func (h *Handlers) GetUser(c echo.Context) error {
	//id, _ := strconv.Atoi(c.Param("id"))
	return c.String(http.StatusOK, "Not implemented")
}

func (h *Handlers) UpdateUser(c echo.Context) error {
	return c.String(http.StatusOK, "Not implemented")
}

func (h *Handlers) DeleteUser(c echo.Context) error {
	return c.String(http.StatusOK, "Not implemented")
}

func (h *Handlers) GetAllUsers(c echo.Context) error {
	return c.String(http.StatusOK, "Not implemented")
}
