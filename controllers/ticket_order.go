package controllers

import (
	"back-end-golang/dtos"
	"back-end-golang/helpers"
	"back-end-golang/middlewares"
	"back-end-golang/usecases"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type TicketOrderController interface {
	GetTicketOrders(c echo.Context) error
	GetTicketOrdersByAdmin(c echo.Context) error
	GetTicketOrderDetailByAdmin(c echo.Context) error
	GetTicketOrderByID(c echo.Context) error
	CreateTicketOrder(c echo.Context) error
	UpdateTicketOrder(c echo.Context) error
}

type ticketOrderController struct {
	ticketOrderUsecase usecases.TicketOrderUsecase
}

func NewTicketOrderController(ticketOrderUsecase usecases.TicketOrderUsecase) TicketOrderController {
	return &ticketOrderController{ticketOrderUsecase}
}

func (c *ticketOrderController) GetTicketOrders(ctx echo.Context) error {
	tokenString := middlewares.GetTokenFromHeader(ctx.Request())
	if tokenString == "" {
		return ctx.JSON(
			http.StatusUnauthorized,
			helpers.NewErrorResponse(
				http.StatusUnauthorized,
				"No token provided",
				helpers.GetErrorData(nil),
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

	pageParam := ctx.QueryParam("page")
	page, err := strconv.Atoi(pageParam)
	if err != nil {
		page = 1
	}

	limitParam := ctx.QueryParam("limit")
	limit, err := strconv.Atoi(limitParam)
	if err != nil {
		limit = 1000
	}

	statusParam := ctx.QueryParam("status")

	ticketOrder, count, err := c.ticketOrderUsecase.GetTicketOrders(page, limit, userId, statusParam)
	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed to get a ticket order",
				helpers.GetErrorData(err),
			),
		)
	}

	return ctx.JSON(
		http.StatusOK,
		helpers.NewPaginationResponse(
			http.StatusOK,
			"Successfully to get order tickets",
			ticketOrder,
			page,
			limit,
			count,
		),
	)
}

func (c *ticketOrderController) GetTicketOrdersByAdmin(ctx echo.Context) error {
	pageParam := ctx.QueryParam("page")
	page, err := strconv.Atoi(pageParam)
	if err != nil {
		page = 1
	}

	limitParam := ctx.QueryParam("limit")
	limit, err := strconv.Atoi(limitParam)
	if err != nil {
		limit = 1000
	}

	searchParam := ctx.QueryParam("search")
	dateStartParam := ctx.QueryParam("date_start")
	dateEndParam := ctx.QueryParam("date_end")
	orderByParam := ctx.QueryParam("order_by")
	filterParam := ctx.QueryParam("filter")

	ticketOrder, count, err := c.ticketOrderUsecase.GetTicketOrdersByAdmin(page, limit, searchParam, dateStartParam, dateEndParam, orderByParam, filterParam)
	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed to get a ticket order",
				helpers.GetErrorData(err),
			),
		)
	}

	return ctx.JSON(
		http.StatusOK,
		helpers.NewPaginationResponse(
			http.StatusOK,
			"Successfully to get order tickets",
			ticketOrder,
			page,
			limit,
			count,
		),
	)
}

func (c *ticketOrderController) GetTicketOrderDetailByAdmin(ctx echo.Context) error {
	trainIdParam := ctx.QueryParam("train_id")
	trainId, _ := strconv.Atoi(trainIdParam)
	ticketOrderIdParam := ctx.QueryParam("ticket_order_id")
	ticketOrderId, _ := strconv.Atoi(ticketOrderIdParam)

	ticketOrder, err := c.ticketOrderUsecase.GetTicketOrdersDetailByAdmin(uint(ticketOrderId), uint(trainId))
	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed to get a ticket order",
				helpers.GetErrorData(err),
			),
		)
	}

	return ctx.JSON(
		http.StatusOK,
		helpers.NewResponse(
			http.StatusOK,
			"Successfully to get order tickets",
			ticketOrder,
		),
	)
}

func (c *ticketOrderController) GetTicketOrderByID(ctx echo.Context) error {
	tokenString := middlewares.GetTokenFromHeader(ctx.Request())
	if tokenString == "" {
		return ctx.JSON(
			http.StatusUnauthorized,
			helpers.NewErrorResponse(
				http.StatusUnauthorized,
				"No token provided",
				helpers.GetErrorData(nil),
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

	trainIdParam := ctx.QueryParam("train_id")
	trainId, _ := strconv.Atoi(trainIdParam)
	ticketOrderIdParam := ctx.QueryParam("ticket_order_id")
	ticketOrderId, _ := strconv.Atoi(ticketOrderIdParam)

	ticketOrder, err := c.ticketOrderUsecase.GetTicketOrderByID(userId, uint(ticketOrderId), uint(trainId))
	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed to get a ticket order",
				helpers.GetErrorData(err),
			),
		)
	}

	return ctx.JSON(
		http.StatusOK,
		helpers.NewResponse(
			http.StatusOK,
			"Successfully to get order tickets",
			ticketOrder,
		),
	)
}

func (c *ticketOrderController) CreateTicketOrder(ctx echo.Context) error {
	tokenString := middlewares.GetTokenFromHeader(ctx.Request())
	if tokenString == "" {
		return ctx.JSON(
			http.StatusUnauthorized,
			helpers.NewErrorResponse(
				http.StatusUnauthorized,
				"No token provided",
				helpers.GetErrorData(nil),
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

	var ticketOrderInput dtos.TicketOrderInput
	if err := ctx.Bind(&ticketOrderInput); err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed binding ticket order",
				helpers.GetErrorData(err),
			),
		)
	}

	ticketOrder, err := c.ticketOrderUsecase.CreateTicketOrder(userId, ticketOrderInput)
	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed to created a ticket order",
				helpers.GetErrorData(err),
			),
		)
	}

	return ctx.JSON(
		http.StatusCreated,
		helpers.NewResponse(
			http.StatusCreated,
			"Successfully to created a ticket order",
			ticketOrder,
		),
	)
}

func (c *ticketOrderController) UpdateTicketOrder(ctx echo.Context) error {
	tokenString := middlewares.GetTokenFromHeader(ctx.Request())
	if tokenString == "" {
		return ctx.JSON(
			http.StatusUnauthorized,
			helpers.NewErrorResponse(
				http.StatusUnauthorized,
				"No token provided",
				helpers.GetErrorData(nil),
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

	var ticketOrderInput dtos.TicketOrderInput
	if err := ctx.Bind(&ticketOrderInput); err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed binding ticket order",
				helpers.GetErrorData(err),
			),
		)
	}

	ticketOrderIDParam := ctx.QueryParam("ticket_order_id")
	ticketOrderID, _ := strconv.Atoi(ticketOrderIDParam)

	statusParam := ctx.QueryParam("status")

	ticketOrder, err := c.ticketOrderUsecase.UpdateTicketOrder(userId, uint(ticketOrderID), statusParam)
	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed to update a ticket order",
				helpers.GetErrorData(err),
			),
		)
	}

	return ctx.JSON(
		http.StatusCreated,
		helpers.NewResponse(
			http.StatusCreated,
			"Successfully to update a ticket order",
			ticketOrder,
		),
	)

}
