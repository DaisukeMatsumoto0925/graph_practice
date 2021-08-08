package server

import (
	"log"
	"os"

	"github.com/labstack/echo"
)

func Run(handler *echo.Echo) {
	handler.Logger.Fatal(handler.Start(":" + os.Getenv("PORT")))
	log.Println("Server exiting")
}
