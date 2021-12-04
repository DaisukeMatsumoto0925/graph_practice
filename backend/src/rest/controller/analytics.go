package controller

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net/http"

	"github.com/DaisukeMatsumoto0925/backend/src/domain"
	"github.com/labstack/echo"
)

type AnalyticsController struct {
	controller *Controller
}

func NewAnalyticsController(controller *Controller) *AnalyticsController {
	return &AnalyticsController{
		controller: controller,
	}
}

func (a *AnalyticsController) Sessions(c echo.Context) {
	var tasks []domain.Task

	if err := a.controller.db.Find(&tasks).Error; err != nil {
		fmt.Println("error")
		return
	}

	fmt.Println(tasks)

	c.Response().Writer.Header().Set("Content-Disposition", "attachment; filename=task.csv")
	c.Response().Writer.Header().Set("Content-Type", "text/csv")

	bin := []byte("\x56\x34\x12\x00\xFF")
	reader := bytes.NewReader(bin)
	binary.Read(reader, binary.LittleEndian, &tasks)
	c.Blob(http.StatusOK, "text/csv", bin)
}
