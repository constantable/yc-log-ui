package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type User struct {
	Id       uint64 `json:"id,omitempty"`
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
	IsAdmin  bool   `json:"isAdmin"`
}

func (h *Handlers) CreateUser(c echo.Context) error {
	userReq := User{}

	if err := c.Bind(&userReq); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if len(userReq.Username) < 3 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Username should be more than 3 chars length",
		})
	}
	if len(userReq.Password) < 3 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Password should be more than 3 chars length",
		})
	}

	user, err := dbHandle.CreateUser(userReq.Username, userReq.Password, userReq.IsAdmin)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, User(user))
}

func (h *Handlers) GetUser(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	user, err := dbHandle.GetUserById(uint64(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}
	return c.JSON(http.StatusOK, User(user))
}

func (h *Handlers) UpdateUser(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	if id == 1 {
		return c.JSON(http.StatusForbidden, map[string]string{
			"message": "Can't update root",
		})
	}
	userReq := User{}
	if err := c.Bind(&userReq); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	user, err := dbHandle.UpdateUser(uint64(id), userReq.Password, userReq.IsAdmin)
	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}
	return c.JSON(http.StatusOK, User(user))
}

func (h *Handlers) DeleteUser(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	if id == 1 {
		return c.JSON(http.StatusForbidden, map[string]string{
			"message": "Can't delete root",
		})
	}
	err := dbHandle.RemoveUser(uint64(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}
	return c.JSON(http.StatusOK, nil)
}

func (h *Handlers) GetAllUsers(c echo.Context) error {
	users, err := dbHandle.GetUsers()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	var ret []User
	for _, user := range users {
		ret = append(ret, User(user))
	}
	return c.JSON(http.StatusOK, ret)
}
