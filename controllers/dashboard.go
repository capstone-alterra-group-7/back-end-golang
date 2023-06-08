package controllers

import (
	"back-end-golang/helpers"
	"back-end-golang/usecases"
	"net/http"

	"github.com/labstack/echo/v4"
)

type DashboardController interface {
	DashboardGetAll(c echo.Context) error
}

type dashboardController struct {
	dashboardUsecase usecases.DashboardUsecase
}

func NewDashboardController(dashboardUsecase usecases.DashboardUsecase) DashboardController {
	return &dashboardController{dashboardUsecase}
}

// Implementasi fungsi-fungsi dari interface ItemController

func (c *dashboardController) DashboardGetAll(ctx echo.Context) error {
	dashboards, err := c.dashboardUsecase.DashboardGetAll()
	if err != nil {
		return ctx.JSON(
			http.StatusInternalServerError,
			helpers.NewErrorResponse(
				http.StatusInternalServerError,
				"Failed fetching dashboard",
				helpers.GetErrorData(err),
			),
		)
	}

	return ctx.JSON(
		http.StatusOK,
		helpers.NewResponse(
			http.StatusOK,
			"Successfully get all dashboards",
			dashboards,
		),
	)
}
