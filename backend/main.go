package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"net/http"
	"os"

	"github.com/constantable/yc-log-ui/access_middleware"
	"github.com/constantable/yc-log-ui/db"
	"github.com/constantable/yc-log-ui/handlers"
	"github.com/joho/godotenv"
)

var dbHandle *db.DB

func init() {
	err := godotenv.Load(".env.local")

	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dbFile := os.Getenv("DBFILE")

	dbHandle, err = db.NewDB(dbFile)
	if err != nil {
		log.Fatal("Error loading database", err)
	}
}

func main() {
	defer dbHandle.Close()

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Welcome to YC Log UI")
	})
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
		AllowCredentials: true,
	}))

	h := &handlers.Handlers{}
	h.SetDatabase(dbHandle)
	e.POST("/api/login", h.Login)
	e.POST("/api/token", h.Token)

	r := e.Group("/api")

	r.Use(access_middleware.IsLoggedIn)
	r.POST("/logs", h.GetLogs)
	r.POST("/logs/page", h.GetPageLogs)
	r.POST("/password", h.Password)

	r.GET("/users", h.GetAllUsers, access_middleware.IsAdmin)
	r.POST("/users", h.CreateUser, access_middleware.IsAdmin)
	r.GET("/users/:id", h.GetUser, access_middleware.IsAdmin)
	r.PUT("/users/:id", h.UpdateUser, access_middleware.IsAdmin)
	r.DELETE("/users/:id", h.DeleteUser, access_middleware.IsAdmin)

	port := os.Getenv("BACKEND_PORT")
	e.Logger.Fatal(e.Start(":" + port))
}
