package controllers

import (
	"back-end-golang/helpers"
	"back-end-golang/middlewares"
	"back-end-golang/usecases"
	"net/http"

	"github.com/labstack/echo/v4"
)

type NotificationController interface {
	GetNotificationByUserID(c echo.Context) error
}

type notificationController struct {
	notificationUsecase usecases.NotificationUsecase
}

func NewNotificationController(notificationUsecase usecases.NotificationUsecase) NotificationController {
	return &notificationController{notificationUsecase}
}

// Implementasi fungsi-fungsi dari interface ItemController

func (c *notificationController) GetNotificationByUserID(ctx echo.Context) error {
	tokenString := middlewares.GetTokenFromHeader(ctx.Request())
	if tokenString == "" {
		return ctx.JSON(
			http.StatusUnauthorized,
			helpers.NewErrorResponse(
				http.StatusUnauthorized,
				"No token provided",
				"Unauthorized",
			),
		)
	}

	userId, err := middlewares.GetUserIdFromToken(tokenString)
	if err != nil {
		return ctx.JSON(
			http.StatusUnauthorized,
			helpers.NewErrorResponse(
				http.StatusUnauthorized,
				"No token provided",
				helpers.GetErrorData(err),
			),
		)
	}
	notification, err := c.notificationUsecase.GetNotificationByUserID(userId)

	if err != nil {
		return ctx.JSON(
			http.StatusNotFound,
			helpers.NewErrorResponse(
				http.StatusNotFound,
				"Failed to get notification by user id",
				helpers.GetErrorData(err),
			),
		)
	}

	return ctx.JSON(
		http.StatusOK,
		helpers.NewResponse(
			http.StatusOK,
			"Successfully to get notification by user id",
			notification,
		),
	)

}
