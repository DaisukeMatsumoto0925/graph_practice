package controller

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"

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

func (a *AnalyticsController) TaskCSV(c echo.Context) error {
	var tasks []domain.Task

	if err := a.controller.db.Find(&tasks).Error; err != nil {
		fmt.Println("error")
		return err
	}

	fmt.Println(tasks)

	csvBytes, err := convertCSV(tasks)
	if err != nil {
		fmt.Println("convertError:", err)
		return err
	}
	newCSVResponse(c, http.StatusOK, csvBytes)

	return nil
}

func newCSVResponse(c echo.Context, status int, data []byte) {
	c.Response().Writer.Header().Set("Content-Disposition", "attachment; filename=task.csv")
	c.Response().Writer.Header().Set("Content-Type", "text/csv")
	c.Blob(status, "text/csv", data)
}

func convertCSV(tasks []domain.Task) ([]byte, error) {
	b := new(bytes.Buffer)
	w := csv.NewWriter(b)

	var header = []string{
		"id",
		"userID",
		"title",
		"note",
		"completed",
		"createdAt",
		"updatedAt",
	}
	w.Write(header)

	for _, task := range tasks {

		var col = []string{
			strconv.Itoa(task.ID),
			strconv.Itoa(task.UserID),
			task.Title,
			task.Note,
			strconv.Itoa(task.Completed),
			task.CreatedAt.String(),
			task.UpdatedAt.String(),
		}
		w.Write(col)
	}
	w.Flush()

	if err := w.Error(); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}
